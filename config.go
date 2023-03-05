package main

import (
	"os"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
 log   "github.com/sirupsen/logrus"
)

type Config struct {
	dav string
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
	dav = viper.GetString("dav")
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		log.Info("! creating directories and files for quickstart")

		mustMkdirAll("./data/public-writable")
		mustWriteFile("./data/public-writable/1.txt", "This is the content of 1.txt")

		mustMkdirAll("./data/public-readonly")
		mustWriteFile("./data/public-readonly/2.txt", "This is the content of 2.txt")

		mustMkdirAll("./data/private-writable")
		mustWriteFile("./data/private-writable/3.txt", "This is the content of 3.txt")
	}

	return
}

func (config *Config) Load() error {

	viper.Set("Verbose", true)
	viper.SetDefault("dav","/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false")
	viper.SetDefault("listen","80"); 
	viper.SetDefault("ip", "0.0.0.0");
	viper.SetDefault("ssl", "no");
	viper.SetDefault("ssl_certificate", "localhost.crt")
	viper.SetDefault("ssl_certificate_key", "localhost.key")
	viper.SetDefault("loglevel", "info")

   	

	pflag.String("dav", "", "like /dav1,./TestDir1,user1,pass1,false;/dav2,./TestDir2,user2,pass2,false")
	pflag.String("listen","80","like 80"); 
	pflag.String("ip", "localhost","like localhost");
	pflag.Bool("ssl",false,"for https ");
	pflag.String("ssl_certificate","localhost.crt", "like localhost.crt")
	pflag.String("ssl_certificate_key","localhost.key", "like localhost.key")
	pflag.String("loglevel","info", "info/debug/error")
	
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Error("! viper.BindPFlags: ", err)
	}
	err = viper.ReadInConfig()
	if err != nil {
		log.Error("! viper.ReadInConfig: ",err)
		config.dav = prepareQuickStart()
	} else {
		config.dav = viper.GetString("dav")
	}

	configurateLogger()	

	log.Debug("dav:", viper.GetString("dav"))
	log.Debug("listen:", viper.GetString("listen"))
	log.Debug("ip:", viper.GetString("ip"))
	log.Debug("ssl:", viper.GetBool("ssl"))
	log.Debug("ssl_certificate:", viper.GetString("ssl_certificate"))
	log.Debug("ssl_certificate_key:", viper.GetString("ssl_certificate_key"))
	log.Debug("loglevel:", viper.GetString("loglevel"))

if viper.GetBool("ssl") == true {
	if _, err := os.Stat(viper.GetString("ssl_certificate")); os.IsNotExist(err) {
		log.Error("ssl_certificate: ", viper.GetString("ssl_certificate"), " NotExist")
		return err
		}
	if _, err := os.Stat(viper.GetString("ssl_certificate_key")); os.IsNotExist(err) {
		log.Error("ssl_certificate_key: ", viper.GetString("ssl_certificate_key"), " NotExist")
		return err
		}

	}
	return nil
}

var AppConfig *Config = &Config{}
