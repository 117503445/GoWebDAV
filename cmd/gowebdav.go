package main

import "GoWebDAV/internal/server"

func main() {
	server := server.NewWebDAVServer("localhost:8080", []*server.HandlerConfig{
		{
			Prefix:   "/data1",
			PathDir:  "./data/public-writable",
			Username: "",
			Password: "",
			ReadOnly: false,
		}, {
			Prefix:   "/data2",
			PathDir:  "./data/public-writable",
			Username: "",
			Password: "",
			ReadOnly: false,
		},
	})
	server.Run()
}
