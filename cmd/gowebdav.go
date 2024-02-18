package main

import (
	_ "GoWebDAV/internal/common"
	"GoWebDAV/internal/server"
)

func main() {
	server, err := server.NewWebDAVServer("localhost:8080", []*server.HandlerConfig{
		{
			Prefix:   "/data1",
			PathDir:  "./data/public-writable",
			Username: "u1",
			Password: "p1",
			ReadOnly: false,
		}, {
			Prefix:   "/data2",
			PathDir:  "./data/public-writable",
			Username: "u2",
			Password: "p2",
			ReadOnly: false,
		},
	})
	if err != nil {
		panic(err)
	}
	server.Run()
}
