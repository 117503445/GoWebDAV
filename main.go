package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/webdav"
)

func main() {
	err := http.ListenAndServe(":8080", &webdav.Handler{

		FileSystem: webdav.Dir("./testdir"),

		LockSystem: webdav.NewMemLS(),
	})
	if err != nil {
		fmt.Println(err)
	}
}
