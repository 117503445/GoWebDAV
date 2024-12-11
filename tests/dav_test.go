package tests

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/117503445/GoWebDAV/internal/server"

	"github.com/stretchr/testify/assert"
	"github.com/studio-b12/gowebdav"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
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

func uiTest(assert *assert.Assertions, url string) {
	resp, err := http.Get(url)
	assert.Nil(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	assert.Nil(err)

	// isHTML := bytes.Equal(body, server.WebdavjsHTML) || bytes.Equal(body, server.WebdavjsHTML_RO)
	isHTML := bytes.Equal(body, server.WebdavjsHTML)
	assert.True(isHTML)
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

	server, err := server.NewWebDAVServer(ADDR, []*server.HandlerConfig{
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
	}, false, "")
	assert.Nil(err)
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

	urls := []string{
		"http://" + ADDR + "/",
		"http://" + ADDR + "/public-writeable",
		"http://" + ADDR + "/public-readonly",
	}
	for _, url := range urls {
		uiTest(assert, url)
	}
}

func TestSingleDav(t *testing.T) {
	const ADDR = "localhost:8082"
	assert := assert.New(t)

	dir := "./data/single-dav"
	os.Mkdir(dir, 0755)

	assert.Nil(os.WriteFile(dir+"/"+file1Name, file1Content, 0644))

	server, err := server.NewWebDAVServer(ADDR, []*server.HandlerConfig{
		{
			Prefix:   "/",
			PathDir:  dir,
			Username: "",
			Password: "",
			ReadOnly: false,
		},
	}, false, "")
	assert.Nil(err)
	go server.Run()

	client := gowebdav.NewClient("http://"+ADDR+"/", "", "")
	untilConnected(assert, client)

	apiTest(assert, client)
}

func TestPlugin(t *testing.T) {
	ast := assert.New(t)
	const src = `package foo
func Bar(s string) string { return s + "-Foo" }`
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)
	_, err := i.Eval(src)
	if err != nil {
		panic(err)
	}
	v, err := i.Eval("foo.Bar")
	if err != nil {
		panic(err)
	}
	bar := v.Interface().(func(string) string)
	r := bar("Kung")
	ast.Equal("Kung-Foo", r)

	// test concurrent
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := bar("Kung")
			ast.Equal("Kung-Foo", r)
		}()
	}
	wg.Wait()

	// https://github.com/traefik/yaegi/discussions/1271
	// custom struct
	i = interp.New(interp.Options{})
	i.Use(stdlib.Symbols)

	type Data struct {
		Message string
	}

	d := &Data{Message: "Kung"}
	custom := make(map[string]map[string]reflect.Value)
	custom["custom/custom"] = make(map[string]reflect.Value)
	custom["custom/custom"]["Data"] = reflect.ValueOf((*Data)(nil))
	i.Use(custom)

	_, err = i.Eval(`
	import "custom"
	func Bar(d *custom.Data) string { return d.Message + "-Foo" }`)
	if err != nil {
		panic(err)
	}

	v, err = i.Eval("Bar")
	if err != nil {
		panic(err)
	}

	newBar := v.Interface().(func(*Data) string)

	newR := newBar(d)
	println(newR)
}

func TestPreRequest(t *testing.T) {
	const ADDR = "localhost:8083"
	assert := assert.New(t)

	dav1Dir := "./data/pre-request-dav1"
	assert.Nil(os.MkdirAll(dav1Dir, 0755))
	assert.Nil(os.WriteFile(dav1Dir+"/"+file1Name, file1Content, 0644))

	dav2Dir := "./data/pre-request-dav2"
	assert.Nil(os.MkdirAll(dav2Dir, 0755))

	assert.Nil(os.WriteFile(dav2Dir+"/"+file1Name, file1Content, 0644))

	server, err := server.NewWebDAVServer(ADDR, []*server.HandlerConfig{
		{
			Prefix:   "/dav1",
			PathDir:  dav1Dir,
			Username: "",
			Password: "",
			ReadOnly: false,
		},
		{
			Prefix:   "/dav2",
			PathDir:  dav2Dir,
			Username: "",
			Password: "",
			ReadOnly: false,
		},
	}, false, "/workspace/assets/Plugins/PreRequestExample.go")
	assert.Nil(err)
	go server.Run()

	// annonymous can't access dav1
	client := gowebdav.NewClient("http://"+ADDR+"/", "", "")
	untilConnected(assert, client)
	_, err = client.Read(fmt.Sprintf("/dav1/%s", file1Name))
	assert.Error(err)

	// user2 can only read dav1
	client = gowebdav.NewClient("http://"+ADDR+"/dav1", "user2", "pass2")
	untilConnected(assert, client)
	readOnlyAPITest(assert, client)

	// user1 can read/write dav1
	client = gowebdav.NewClient("http://"+ADDR+"/dav1", "user1", "pass1")
	untilConnected(assert, client)
	apiTest(assert, client)

	// user3 can't access dav1
	client = gowebdav.NewClient("http://"+ADDR+"/dav1", "user3", "pass3")
	untilConnected(assert, client)
	_, err = client.Read(fmt.Sprintf("/dav1/%s", file1Name))
	assert.Error(err)

	// everyone access dav2
	client = gowebdav.NewClient("http://"+ADDR+"/dav2", "", "")
	untilConnected(assert, client)
	apiTest(assert, client)
}
