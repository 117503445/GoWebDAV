package tests

import (
	"os"
	"testing"
	"time"

	"GoWebDAV/internal/server"
	"github.com/stretchr/testify/assert"
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

var file1Name = "1.txt"
var file1Content = []byte("file1")

var file2Name = "2.txt"
var file2Content = []byte("file2")

func untilConnected(assert *assert.Assertions, client *gowebdav.Client) {
	maxTries := 10
	for {
		err := client.Connect()
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
		maxTries--
		if maxTries == 0 {
			assert.FailNow("", "Failed to connect to the server")
		}
	}
}

// readOnlyAPITest tests the following:
// 1. Read the file
// 2. Read the directory
// before test, the client should be connected to the server and one file should be created
func readOnlyAPITest(assert *assert.Assertions, client *gowebdav.Client) {
	downloadContent, err := client.Read("/1.txt")
	assert.Nil(err)

	if string(downloadContent) != string(file1Content) {
		assert.FailNow("Failed to download the file")
	}

	files, err := client.ReadDir("/")
	assert.Nil(err)
	if len(files) != 1 {
		assert.FailNowf("files number not equal to 1", "len(files) != 1, got %d", len(files))
	}
}

// apiTest tests the following:
// 1. Create a directory
// 2. Write a file to the directory
// 3. Read the directory
// 4. Read the file
// 5. Remove the file
// before test, the client should be connected to the server and one file should be created
func apiTest(assert *assert.Assertions, client *gowebdav.Client) {
	readOnlyAPITest(assert, client)

	assert.Nil(client.Mkdir("/dir", os.ModePerm))
	client.Write("/dir/"+file2Name, file2Content, 0644)

	files, err := client.ReadDir("/dir")
	assert.Nil(err)
	if len(files) != 1 {
		assert.FailNowf("files number not equal to 1", "len(files) != 1, got %d", len(files))
	}

	downloadContent, err := client.Read("/dir/2.txt")
	assert.Nil(err)
	if string(downloadContent) != string(file2Content) {
		assert.FailNow("Failed to download the file")
	}

	assert.Nil(client.Remove("/dir/2.txt"))

	files, err = client.ReadDir("/dir")
	assert.Nil(err)
	if len(files) != 0 {
		assert.FailNowf("files number not equal to 0", "len(files) != 0, got %d", len(files))
	}
}

func TestMultiDav(t *testing.T) {
	const ADDR = "localhost:8081"
	assert := assert.New(t)

	dirPublicWriteable := "./data/public-writeable"
	dirPublicReadonly := "./data/public-readonly"
	dirAuthWriteable := "./data/auth-writeable"

	dirs := []string{dirPublicWriteable, dirPublicReadonly, dirAuthWriteable}
	for _, dir := range dirs {
		assert.Nil(os.Mkdir(dir, 0755))
		assert.Nil(os.WriteFile(dir+"/"+file1Name, file1Content, 0644))
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

	untilConnected(assert, publicWriteableClient)
	untilConnected(assert, publicReadonlyClient)
	untilConnected(assert, authWriteableClient)

	apiTest(assert, publicWriteableClient)
	readOnlyAPITest(assert, publicReadonlyClient)
	apiTest(assert, authWriteableClient)

	indexClient := gowebdav.NewClient("http://"+ADDR+"/", "", "")
	untilConnected(assert, indexClient)

	files, err := indexClient.ReadDir("/")
	assert.Nil(err)
	if len(files) != 3 {
		assert.FailNowf("files number not equal to 3", "len(files) != 3, got %d", len(files))
	}
}

func TestSingleDav(t *testing.T) {
	const ADDR = "localhost:8082"
	assert := assert.New(t)

	dir := "./data/single-dav"
	os.Mkdir(dir, 0755)

	assert.Nil(os.WriteFile(dir+"/"+file1Name, file1Content, 0644))

	server := server.NewWebDAVServer(ADDR, []*server.HandlerConfig{
		{
			Prefix:   "/",
			PathDir:  dir,
			Username: "",
			Password: "",
			ReadOnly: false,
		},
	})
	go server.Run()

	client := gowebdav.NewClient("http://"+ADDR+"/", "", "")
	untilConnected(assert, client)

	apiTest(assert, client)
}
