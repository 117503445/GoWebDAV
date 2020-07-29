package main

import (
	"GoWebdav/model"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func handleDirList(fs webdav.FileSystem, w http.ResponseWriter, req *http.Request, prefix string) bool {
	ctx := context.Background()

	path := req.URL.Path
	path = strings.Replace(path, prefix, "/", 1)

	f, err := fs.OpenFile(ctx, path, os.O_RDONLY, 0)
	if err != nil {
		return false
	}
	defer f.Close()

	if fi, _ := f.Stat(); fi != nil && !fi.IsDir() {

		return false

	}

	dirs, err := f.Readdir(-1)

	if err != nil {
		log.Print(w, "Error reading directory", http.StatusInternalServerError)
		return false
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, err = fmt.Fprintf(w, "<pre>\n")
	if err != nil {
		fmt.Println(err)
	}

	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() {
			name += "/"
		}
		_, err = fmt.Fprintf(w, "<a href=\"%s\" >%s</a>\n", name, name)
		if err != nil {
			fmt.Println(err)
		}
	}

	_, err = fmt.Fprintf(w, "</pre>\n")
	if err != nil {
		fmt.Println(err)
	}
	return true

}
func main() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	AppConfig.Load()

	davConfigs := strings.Split(AppConfig.dav, ";")

	WebDAVConfigs := make([]*model.WebDAVConfig, 0)

	for _, davConfig := range davConfigs {
		davConfigArray := strings.Split(davConfig, ",")
		prefix := davConfigArray[0]
		pathDir := davConfigArray[1]
		username := davConfigArray[2]
		password := davConfigArray[3]

		WebDAVConfig := &model.WebDAVConfig{}
		WebDAVConfig.Init(prefix, pathDir, username, password)
		WebDAVConfigs = append(WebDAVConfigs, WebDAVConfig)
	}

	sMux := http.NewServeMux()
	sMux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {

		username, password, ok := req.BasicAuth()

		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		webDAVConfig := model.WebDAVConfigFindOneByPrefix(WebDAVConfigs, parsePrefixFromURL(req.URL))
		if webDAVConfig == nil {
			http.NotFound(w, req)
			return
		}

		if username != webDAVConfig.Username || password != webDAVConfig.Password {
			http.Error(w, "WebDAV: need authorized!", http.StatusUnauthorized)
			return
		}

		if req.Method == "GET" && handleDirList(webDAVConfig.Handler.FileSystem, w, req, webDAVConfig.Handler.Prefix) {
			return
		}

		webDAVConfig.Handler.ServeHTTP(w, req)

	})

	err := http.ListenAndServe(":8080", sMux)
	if err != nil {
		fmt.Println(err)
	}
}

// /dav1/123.txt -> dav1
func parsePrefixFromURL(url *url.URL) string {
	u := fmt.Sprint(url)
	return "/" + strings.Split(u, "/")[1]
}
