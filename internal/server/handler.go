package server

import (
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/webdav"
)

//go:embed webdavjs.html
var WebdavjsHTML string

type HandlerConfig struct {
	Prefix   string
	PathDir  string
	Username string
	Password string
	ReadOnly bool
}

type handler struct {
	handler *webdav.Handler

	prefix  string // URL prefix
	dirPath string // File system directory

	username string // HTTP Basic Auth Username. if empty, no auth
	password string // HTTP Basic Auth Password.

	readOnly bool // if true, only allow GET, OPTIONS, PROPFIND, HEAD
}

func NewHandler(cfg *HandlerConfig) *handler {
	return &handler{
		handler: &webdav.Handler{
			FileSystem: webdav.Dir(cfg.PathDir),
			LockSystem: webdav.NewMemLS(),
			Prefix:     cfg.Prefix,
		},
		prefix:   cfg.Prefix,
		dirPath:  cfg.PathDir,
		username: cfg.Username,
		password: cfg.Password,
		readOnly: cfg.ReadOnly,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	enableBasicAuth := h.username != "" && h.username != "null"
	if enableBasicAuth {
		username, password, ok := req.BasicAuth()
		// log.Debug().Str("username", username).Str("password", password).Bool("ok", ok).Msg("BasicAuth Request")
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if username != h.username || password != h.password {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	if h.readOnly {
		allowedMethods := map[string]bool{
			"GET":      true,
			"OPTIONS":  true,
			"PROPFIND": true,
			"HEAD":     true,
		}
		if !allowedMethods[req.Method] {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}

	log.Debug().Str("URL", req.URL.Path).Str("Method", req.Method).Msg("Handler Request")
	if req.Method == "GET" && (req.URL.Path == h.prefix || strings.HasSuffix(req.URL.Path, "/")) {
		if _, err := w.Write([]byte(WebdavjsHTML)); err != nil {
			log.Error().Err(err).Msg("Failed to write index.html")
		}
		return
	}

	if req.Method == "HEAD" {
		// for ui
		return
	}

	h.handler.ServeHTTP(w, req)
}

// checkHandlerConfig checks if the handler config is valid
// if mkdir is true, it will create the directory if not exist
func checkHandlerConfig(cfg *HandlerConfig, mkdir bool) error {
	// prefix must start with "/", contains only one "/"
	if cfg.Prefix == "" || cfg.Prefix[0] != '/' {
		return errors.New("prefix must start with /")
	}

	// prefix must contain only one "/"
	if strings.Count(cfg.Prefix, "/") != 1 {
		return errors.New("prefix must contain only one /")
	}

	// prefix must not contain not allowed characters
	notAllowedChars := []string{"?", "%", "#", "&"}
	for _, char := range notAllowedChars {
		if strings.Contains(cfg.Prefix, char) {
			return errors.New("prefix must not contain " + char)
		}
	}

	// pathDir must be a valid directory
	if fileinfo, err := os.Stat(cfg.PathDir); os.IsNotExist(err) {
		if !mkdir {
			return errors.New("pathDir does not exist")
		} else {
			// try to create the directory
			log.Info().Str("path", cfg.PathDir).Msg("Creating dir")
			if err := os.MkdirAll(cfg.PathDir, 0755); err != nil {
				return err
			}
		}
	} else if err != nil {
		return err
	} else if !fileinfo.IsDir() {
		return errors.New("pathDir must be a directory")
	}

	if cfg.Username != "" && cfg.Password == "" {
		return errors.New("password must not be empty if username is not empty")
	}

	return nil
}

func checkHandlerConfigs(cfgs []*HandlerConfig) error {
	for _, cfg := range cfgs {
		if err := checkHandlerConfig(cfg, true); err != nil {
			return fmt.Errorf("config %+v is invalid: %s", cfg, err.Error())
		}
	}

	if len(cfgs) > 1 {
		for _, cfg := range cfgs {
			if cfg.Prefix == "/" {
				return errors.New("prefix / is not allowed if there are more than one handler")
			}
		}
	}

	// check if prefix is duplicated
	prefixs := make(map[string]bool)
	for _, cfg := range cfgs {
		if _, ok := prefixs[cfg.Prefix]; ok {
			return fmt.Errorf("prefix %s is duplicated", cfg.Prefix)
		}
		prefixs[cfg.Prefix] = true
	}

	return nil
}
