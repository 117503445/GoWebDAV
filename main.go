package main

import (
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/webdav"
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

	fmt.Fprintf(w, "<pre>\n")

	for _, d := range dirs {

		name := d.Name()

		if d.IsDir() {

			name += "/"

		}

		fmt.Fprintf(w, "<a href=\"%s\" >%s</a>\n", name, name)

	}

	fmt.Fprintf(w, "</pre>\n")

	return true

}
func main() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	AppConfig.Load()

	fmt.Println(AppConfig.dav)
	fs := &webdav.Handler{

		FileSystem: webdav.Dir("./TestDir"),

		LockSystem: webdav.NewMemLS(),
		Prefix:     "/dav1",
	}
	sMux := http.NewServeMux()
	sMux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {

		username, password, ok := req.BasicAuth()

		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if username != "user" || password != "123456" {
			http.Error(w, "WebDAV: need authorized!", http.StatusUnauthorized)
			return
		}

		if req.Method == "GET" && handleDirList(fs.FileSystem, w, req, "/dav1") {
			return
		}

		fs.ServeHTTP(w, req)

	})

	err := http.ListenAndServe(":8080", sMux)
	if err != nil {
		fmt.Println(err)
	}
}
