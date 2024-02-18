package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/webdav"

	"GoWebDAV/model"

	"github.com/spf13/viper"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed static/index.html
var indexHTML string

func setLogger() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	log.Logger = log.Output(output)
}

func init() {
	setLogger()
}

func run() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	AppConfig.Load()
	fmt.Printf("dav = %s\n", AppConfig.dav)
	davConfigs := strings.Split(AppConfig.dav, ";")

	WebDAVConfigs := make([]*model.WebDAVConfig, 0)

	for _, davConfig := range davConfigs {
		if len(davConfig) == 0 {
			continue
		}
		WebDAVConfig := model.InitByConfigStr(davConfig)
		// Check for collision
		found, _ := model.ParseURL(WebDAVConfigs, WebDAVConfig.Prefix)
		if found != nil {
			fmt.Printf("Dav names collision: `%s` starts with `%s`", WebDAVConfig.Prefix, found.Prefix)
			os.Exit(1)
		}

		WebDAVConfigs = append(WebDAVConfigs, &WebDAVConfig)
	}

	sMux := http.NewServeMux()
	sMux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {

		webDAVConfig, davPath := model.ParseURL(WebDAVConfigs, req.URL.Path)

		if webDAVConfig == nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")

			// index

			_, err := fmt.Fprintf(w, "<pre>\n")
			if err != nil {
				fmt.Println(err)
			}

			for _, config := range WebDAVConfigs {
				_, err = fmt.Fprintf(w, "<a href=\"%s\" >%s</a>\n", config.Prefix+"/", config.Prefix)
				if err != nil {
					fmt.Println(err)
				}
			}

			_, err = fmt.Fprintf(w, "<pre>\n")
			if err != nil {
				fmt.Println(err)
			}

			return
		}

		// When the username and password in the configuration are both null, no identity check is performed
		if webDAVConfig.Username != "null" && webDAVConfig.Password != "null" {
			username, password, ok := req.BasicAuth()

			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if username == "" || password == "" {
				http.Error(w, "username missing or password missing", http.StatusUnauthorized)
				return
			}

			if username != webDAVConfig.Username || password != webDAVConfig.Password {
				http.Error(w, "username wrong or password wrong", http.StatusUnauthorized)
				return
			}
		}

		if webDAVConfig.ReadOnly {
			allowedMethods := map[string]bool{
				"GET":      true,
				"OPTIONS":  true,
				"PROPFIND": true,
				"HEAD":     true,
			}
			if !allowedMethods[req.Method] {
				w.WriteHeader(http.StatusMethodNotAllowed)
				_, err := w.Write([]byte("Readonly, Method " + req.Method + " Not Allowed"))
				if err != nil {
					fmt.Println(err)
					return
				}
				return
			}
		}

		if req.Method == "GET" && isDir(webDAVConfig.Handler.FileSystem, davPath) {
			_, err := w.Write([]byte(indexHTML))
			if err != nil {
				fmt.Println(err)
				return
			}
			return
		}

		if req.Method == "HEAD" {
			return
		}

		// handle file
		webDAVConfig.Handler.ServeHTTP(w, req)
	})

	addr := fmt.Sprintf("%s:%d", AppConfig.addr, AppConfig.port)

	fmt.Printf("start listen on http://%s\n", addr)
	err := http.ListenAndServe(addr, sMux)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	log.Debug().Msg("Hello")
	run()
}

func isDir(fs webdav.FileSystem, davPath string) bool {
	ctx := context.Background()

	f, err := fs.OpenFile(ctx, davPath, os.O_RDONLY, 0)
	if err != nil {
		return false
	}
	defer f.Close()

	if fi, _ := f.Stat(); fi != nil && !fi.IsDir() {
		return false
	}
	return true
}
