package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	dav string
	port int
}

func mustMkdirAll(path string) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func mustWriteFile(path string, content string) {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
}

// prepareQuickStart create dirs and files for quickstart, and return quickstart dav
func prepareQuickStart() (dav string) {
	dav = "/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false"

	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		fmt.Println("creating directories and files for quickstart")

		mustMkdirAll("./data/public-writable")
		mustWriteFile("./data/public-writable/1.txt", "This is the content of 1.txt")

		mustMkdirAll("./data/public-readonly")
		mustWriteFile("./data/public-readonly/2.txt", "This is the content of 2.txt")

		mustMkdirAll("./data/private-writable")
		mustWriteFile("./data/private-writable/3.txt", "This is the content of 3.txt")
	}

	return
}

func (config *Config) Load() {
	pflag.String("dav", "", "like /dav1,./TestDir1,user1,pass1,false;/dav2,./TestDir2,user2,pass2,false")
	pflag.Int("port", 80, "port")
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		fmt.Println(err)
	}
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	config.dav = viper.GetString("dav")
	config.port = viper.GetInt("port")
	if config.dav == "" {
		config.dav = prepareQuickStart()
	}
}

var AppConfig *Config = &Config{}
