package server

import (
	"context"
	_ "embed"
	"net/http"
	"os"
	"strings"

	"github.com/Jipok/webdavWithPATCH"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/webdav"
)

type WebDAVServer struct {
	addr string // Address to listen on, e.g. "0.0.0.0:8080"
	smux *http.ServeMux
}

func NewWebDAVServer(addr string, handlerConfigs []*HandlerConfig, davListIsSecret bool) (*WebDAVServer, error) {
	if err := checkHandlerConfigs(handlerConfigs); err != nil {
		return nil, err
	}

	sMux := http.NewServeMux()

	handlers := make(map[string]*handler) // URL prefix -> Handler
	for _, cfg := range handlerConfigs {
		h := NewHandler(cfg)
		handlers[cfg.Prefix] = h
	}

	// create a webdav.Handler for listing all available prefixes
	memFileSystem := webdav.NewMemFS()
	for _, cfg := range handlerConfigs {
		if cfg.Prefix == "/" {
			continue
		}
		if err := memFileSystem.Mkdir(context.TODO(), cfg.Prefix, os.ModeDir); err != nil {
			log.Error().Err(err).Str("prefix", cfg.Prefix).Msg("Failed to create directory in memFileSystem")
			continue
		}
	}

	indexHandler := &webdavWithPATCH.Handler{
		Handler: webdav.Handler{
			FileSystem: memFileSystem,
			LockSystem: webdav.NewMemLS(),
			Prefix:     "/",
		},
	}

	sMux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		url := req.URL.Path
		method := req.Method
		log.Debug().Str("URL", url).Str("Method", method).Msg("Request")

		if req.Method == "HEAD" && strings.HasSuffix(req.URL.Path, "/") {
			return // Fixes error for ui
		}

		for prefix, handler := range handlers {
			if !strings.HasPrefix(url, prefix) {
				continue
			}
			handler.ServeHTTP(w, req)
			return
		}

		if davListIsSecret {
			return
		}

		if !readOnlyMethods[req.Method] {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

		if method == "GET" && url == "/" {
			if _, err := w.Write(WebdavjsHTML); err != nil {
				log.Error().Err(err).Msg("Failed to write index.html")
			}
			return
		}

		indexHandler.ServeHTTP(w, req)
	})
	return &WebDAVServer{
		addr: addr,
		smux: sMux,
	}, nil
}

func (s *WebDAVServer) Run() error {
	log.Info().Str("addr", "http://"+s.addr).Msg("WebDAV server started")
	if err := http.ListenAndServe(s.addr, s.smux); err != nil {
		return err
	}
	return nil
}
