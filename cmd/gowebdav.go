package main

import "GoWebDAV/internal/server"

func main() {
	server := server.NewWebDAVServer("localhost:8080", []*server.HandlerConfig{
		{
			Prefix:   "/data",
			PathDir:  "./data",
			Username: "user",
			Password: "pass",
			ReadOnly: false,
		},
	})
	server.Run()
}
