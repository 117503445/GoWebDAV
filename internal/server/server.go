package server

import (
	"context"
	_ "embed"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/webdav"
)

//go:embed index.html
var indexHTML string

type WebDAVServer struct {
	addr string // Address to listen on, e.g. "0.0.0.0:8080"
	smux *http.ServeMux
}

func NewWebDAVServer(addr string, handlerConfigs []*HandlerConfig) *WebDAVServer {
	sMux := http.NewServeMux()

	handlers := make(map[string]*handler) // URL prefix -> Handler
	for _, cfg := range handlerConfigs {
		h := NewHandler(cfg)
		handlers[cfg.Prefix] = h
	}

	// single dav mode: if there is only one handler and its prefix is "/", route all requests to it
	enableSingleDavMode := false
	var singleHandler *handler
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
		url := req.URL.Path
		method := req.Method
		log.Debug().Str("URL", url).Str("Method", method).Msg("Request")

		if enableSingleDavMode {
			// TODO: ui for single dav mode
			singleHandler.ServeHTTP(w, req)
			return
		}

		for prefix, handler := range handlers {
			if !strings.HasPrefix(url, prefix) {
				continue
			}
			if method == "HEAD" && url == prefix+"/" {
				continue
			}
			handler.ServeHTTP(w, req)
			return
		}

		if method == "PROPFIND" && url == "/" {
			indexHandler.ServeHTTP(w, req)
			return
		}

		if method == "GET" && url == "/" {
			// TODO: return dav selector
		}

		log.Debug().Str("URL", url).Str("Method", method).Msg("return html")
		if _, err := w.Write([]byte(indexHTML)); err != nil {
			log.Error().Err(err).Msg("Failed to write index.html")
		}
	})
	return &WebDAVServer{
		addr: addr,
		smux: sMux,
	}
}

func (s *WebDAVServer) Run() {
	if err := http.ListenAndServe(s.addr, s.smux); err != nil {
		panic(err)
	}
}
