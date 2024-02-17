package tests

import (
	"os"
	"testing"
	"time"

	"GoWebDAV/internal/server"
	// "github.com/rs/zerolog/log"
	"github.com/studio-b12/gowebdav"
)

func TestMain(m *testing.M) {
	dirData := "./data"
	if _, err := os.Stat(dirData); err == nil {
		os.RemoveAll(dirData)
	}
	os.Mkdir(dirData, 0755)

	os.Exit(m.Run())
}

const ADDR = "localhost:8081"

func TestPublicWrite(t *testing.T) {
	dirData := "./data/public-writeable"
	os.Mkdir(dirData, 0755)
}

func TestDAV(t *testing.T) {
	var err error
	// Create a directory for the server

	dirPublicWriteable := "./data/public-writeable"
	dirPublicReadonly := "./data/public-readonly"
	dirAuthWriteable := "./data/auth-writeable"

	file1Name := "1.txt"
	file1Content := []byte("file1")

	dirs := []string{dirPublicWriteable, dirPublicReadonly, dirAuthWriteable}
	for _, dir := range dirs {
		if err = os.Mkdir(dir, 0755); err != nil {
			t.Fatal(err)
		}

		if err = os.WriteFile(dir+"/"+file1Name, file1Content, 0644); err != nil {
			t.Fatal(err)
		}
	}

	server := server.NewWebDAVServer(ADDR, []*server.HandlerConfig{
		{
			Prefix:   "/public-writeable",
			PathDir:  dirPublicWriteable,
			Username: "",
			Password: "",
			ReadOnly: false,
		},
		{
			Prefix:   "/public-readonly",
			PathDir:  dirPublicReadonly,
			Username: "",
			Password: "",
			ReadOnly: true,
		},
		{
			Prefix:   "/auth-writeable",
			PathDir:  dirAuthWriteable,
			Username: "user",
			Password: "pass",
			ReadOnly: false,
		},
	})
	go server.Run()

	publicWriteableClient := gowebdav.NewClient("http://"+ADDR+"/public-writeable", "", "")
	publicReadonlyClient := gowebdav.NewClient("http://"+ADDR+"/public-readonly", "", "")
	authWriteableClient := gowebdav.NewClient("http://"+ADDR+"/auth-writeable", "user", "pass")

	untilConnected := func(t *testing.T, client *gowebdav.Client) {
		maxTries := 10
		for {
			err := client.Connect()
			if err == nil {
				break
			}
			time.Sleep(100 * time.Millisecond)
			maxTries--
			if maxTries == 0 {
				t.Fatal("Failed to connect to the server")
			}
		}
	}

	untilConnected(t, publicWriteableClient)
	untilConnected(t, publicReadonlyClient)
	untilConnected(t, authWriteableClient)

	readOnlyAPITest := func(t *testing.T, client *gowebdav.Client) {
		downloadContent, err := client.Read("/1.txt")
		if err != nil {
			t.Fatal(err)
		}
		if string(downloadContent) != string(file1Content) {
			t.Fatal("Failed to download the file")
		}

		files, err := client.ReadDir("/")
		if err != nil {
			t.Fatal(err)
		}
		if len(files) != 1 {
			t.Fatal("len(files) != 1, got", len(files))
		}

	}

	apiTest := func(t *testing.T, client *gowebdav.Client) {
		readOnlyAPITest(t, client)
	}

	apiTest(t, publicWriteableClient)
	readOnlyAPITest(t, publicReadonlyClient)
	apiTest(t, authWriteableClient)
}
