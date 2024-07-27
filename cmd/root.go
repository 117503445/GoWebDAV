package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "GoWebDAV/internal/common"
	"GoWebDAV/internal/server"

	"github.com/117503445/goutils"
	"github.com/rs/zerolog/log"
)

// createDefaultDirs creates default directories when dav is not set
func createDefaultDirs() {
	dirs := []string{
		"./data/public-writable",
		"./data/public-readonly",
		"./data/private-writable",
	}
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				panic(err)
			}

			if err = os.WriteFile(dir+"/1.txt", []byte("Hello"), 0644); err != nil {
				panic(err)
			}
		} else if err != nil {
			panic(err)
		}
	}

}

func parseDavToHandlerConfigs(dav string) (handlerConfigs []*server.HandlerConfig, err error) {
	// dav = "/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false"

	davConfigs := strings.Split(dav, ";")
	for _, davConfig := range davConfigs {
		davConfig = strings.Trim(davConfig, " ")
		if len(davConfig) == 0 {
			continue
		}
		arr := strings.Split(davConfig, ",")
		if len(arr) != 5 && len(arr) != 2 {
			err = fmt.Errorf("invalid dav config: %s", davConfig)
			return
		}
		prefix := arr[0]
		pathDir := arr[1]
		username := ""
		password := ""
		readonly := true
		if len(arr) == 5 {
			username = arr[2]
			password = arr[3]
			readonly, err = strconv.ParseBool(arr[4])
			if err != nil {
				log.Err(err).Msg("Failed to parse readonly")
			}
		}
		handlerConfigs = append(handlerConfigs, &server.HandlerConfig{
			Prefix:   prefix,
			PathDir:  pathDir,
			Username: username,
			Password: password,
			ReadOnly: readonly,
		})
	}
	return handlerConfigs, nil
}

func Execute() {
	const DEFAULT_DAV_CONFIG = "/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false"

	type Config struct {
		Address         string `koanf:"address"`
		Port            string `koanf:"port"`
		Dav             string `koanf:"dav" usage:"dav config, format: prefix,pathDir,username,password,readonly;..."`
		DavListIsSecret bool   `koanf:"secret_dav_list" usage:"if true, hide the dav list"`
	}
	cfg := &Config{
		Address:         "0.0.0.0",
		Port:            "80",
		Dav:             DEFAULT_DAV_CONFIG,
		DavListIsSecret: false,
	}

	goutils.LoadConfig(cfg)

	log.Info().Interface("config", cfg).Msg("Config")

	if cfg.Dav == DEFAULT_DAV_CONFIG {
		createDefaultDirs()
		log.Info().Msg("Default dav config is used, created default directories")
	}

	handlerConfigs, err := parseDavToHandlerConfigs(cfg.Dav)
	if err != nil {
		panic(err)
	}

	for _, handlerConfig := range handlerConfigs {
		log.Debug().Str("prefix", handlerConfig.Prefix).Str("pathDir", handlerConfig.PathDir).Str("username", handlerConfig.Username).Str("password", handlerConfig.Password).Bool("readonly", handlerConfig.ReadOnly).Msg("Dav")
	}

	server, err := server.NewWebDAVServer(cfg.Address+":"+cfg.Port, handlerConfigs, cfg.DavListIsSecret)
	if err != nil {
		panic(err)
	}
	server.Run()
}
