package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type config struct {
	dav string
}

func (config *config) Load() {
	pflag.String("dav", "/dav1,user1,pass1;/dav2,user2,pass2", "like /dav1,user1,pass1;/dav2,user2,pass2")
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

var AppConfig *config = &config{}
