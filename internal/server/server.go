package server

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/webdav"
)

type HandlerConfig struct {
	Prefix   string
	PathDir  string
	Username string
	Password string
	ReadOnly bool
}

type Handler struct {
	handler *webdav.Handler

	prefix  string // URL prefix
	dirPath string // File system directory

	username string // HTTP Basic Auth Username. if empty, no auth
	password string // HTTP Basic Auth Password

	readOnly bool // if true, only allow GET, OPTIONS, PROPFIND, HEAD
}

// func NewHandler(prefix, dirPath, username, password string, readOnly bool) *Handler {
// 	return &Handler{
// 		handler: &webdav.Handler{
// 			FileSystem: webdav.Dir(dirPath),
// 			LockSystem: webdav.NewMemLS(),
// 			Prefix:     prefix,
// 		},
// 	}
// }

func NewHandler(cfg *HandlerConfig) *Handler {
	return &Handler{
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

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	enableBasicAuth := h.username != ""
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

	h.handler.ServeHTTP(w, req)
}

type WebDAVServer struct {
	addr string // Address to listen on, e.g. "0.0.0.0:8080"
	smux *http.ServeMux
}

func (s *WebDAVServer) Run() {
	if err := http.ListenAndServe(s.addr, s.smux); err != nil {
		panic(err)
	}
}

func NewWebDAVServer(addr string, handlerConfigs []*HandlerConfig) *WebDAVServer {
	sMux := http.NewServeMux()

	handlers := make(map[string]*Handler) // URL prefix -> Handler
	for _, cfg := range handlerConfigs {
		h := NewHandler(cfg)
		handlers[cfg.Prefix] = h
	}

	enableSingleDavMode := false
	var singleHandler *Handler
	if len(handlers) == 1 {
		for _, h := range handlers {
			if h.prefix == "/" {
				log.Debug().Msg("Enable SingleDavMode")

				enableSingleDavMode = true
				singleHandler = h
				break
			}
		}
	}

	// create a webdav.Handler for listing all available prefixes
	memFileSystem := webdav.NewMemFS()
	for _, cfg := range handlerConfigs {
		memFileSystem.Mkdir(context.TODO(), cfg.Prefix, os.ModeDir)
	}
	indexHandler := &webdav.Handler{
		FileSystem: memFileSystem,
		LockSystem: webdav.NewMemLS(),
		Prefix:     "/",
	}

	sMux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Debug().Str("URL", req.URL.Path).Str("Method", req.Method).Msg("Request")
		if enableSingleDavMode {
			singleHandler.ServeHTTP(w, req)
			return
		}

		url := req.URL.Path
		for prefix, handler := range handlers {
			if !strings.HasPrefix(url, prefix) {
				continue
			}
			handler.ServeHTTP(w, req)
			return
		}

		if req.Method == "PROPFIND" && url == "/" {
			indexHandler.ServeHTTP(w, req)
			return
		}
	})
	return &WebDAVServer{
		addr: addr,
		smux: sMux,
	}
}
