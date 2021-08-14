package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	dav string
}

func (config *Config) Load() {
	pflag.String("dav", "/dav1,./TestDir1,user1,pass1,false;/dav2,./TestDir2,user2,pass2,false", "like /dav1,./TestDir1,user1,pass1,false;/dav2,./TestDir2,user2,pass2,false")
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
}

var AppConfig *Config = &Config{}
