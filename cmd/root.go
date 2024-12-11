package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/117503445/GoWebDAV/internal/common"
	"github.com/117503445/GoWebDAV/internal/server"

	"github.com/117503445/goutils"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
)

const DEFAULT_DAV_CONFIG = "/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false"

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
				log.Fatal().Err(err).Msg("Failed to create default directories")
			}

			if err = os.WriteFile(dir+"/1.txt", []byte("Hello"), 0644); err != nil {
				log.Fatal().Err(err).Msg("Failed to write file")
			}
		} else if err != nil {
			log.Fatal().Err(err).Msg("Failed to check directory")
		}
	}
}

func isDefaultHandlerConfigs(handlerConfig []*server.HandlerConfig) bool {
	defaultHandlerConfigs, err := parseDavToHandlerConfigs(DEFAULT_DAV_CONFIG)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse default dav")
	}

	if len(handlerConfig) != len(defaultHandlerConfigs) {
		return false
	}

	for i, hc := range handlerConfig {
		if hc.Prefix != defaultHandlerConfigs[i].Prefix ||
			hc.PathDir != defaultHandlerConfigs[i].PathDir ||
			hc.Username != defaultHandlerConfigs[i].Username ||
			hc.Password != defaultHandlerConfigs[i].Password ||
			hc.ReadOnly != defaultHandlerConfigs[i].ReadOnly {
			println(hc.Prefix, defaultHandlerConfigs[i].Prefix)
			return false
		}
	}
	return true
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

func parseDavsToHandlerConfigs(davs []*koanf.Koanf) (handlerConfigs []*server.HandlerConfig, err error) {
	for _, d := range davs {
		prefix := d.String("prefix")
		pathDir := d.String("pathDir")
		username := d.String("username")
		password := d.String("password")
		readOnly := d.Bool("readOnly")
		handlerConfigs = append(handlerConfigs, &server.HandlerConfig{
			Prefix:   prefix,
			PathDir:  pathDir,
			Username: username,
			Password: password,
			ReadOnly: readOnly,
		})
	}
	return handlerConfigs, nil
}

// dav: from cli, config file, env or default. format: prefix,pathDir,username,password,readonly;...
// davs: from config file, which is designed for convenience writing in config file.
func getHandlerConfigs(dav string, davs []*koanf.Koanf) (handlerConfigs []*server.HandlerConfig) {
	// if davs is set and dav == DEFAULT_DAV_CONFIG, use davs and not create default dirs
	// else, merge davs and dav

	handlerConfigs = make([]*server.HandlerConfig, 0)
	var err error

	if len(davs) > 0 {
		davsHandlerConfigs, err := parseDavsToHandlerConfigs(davs)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to parse davs")
		}
		for _, handlerConfig := range davsHandlerConfigs {
			log.Debug().Interface("handlerConfig", handlerConfig).Msg("davs")
		}

		if dav == DEFAULT_DAV_CONFIG {
			return davsHandlerConfigs
		}

		handlerConfigs = append(handlerConfigs, davsHandlerConfigs...)
	}

	davHandlerConfigs, err := parseDavToHandlerConfigs(dav)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse dav")
	}

	for _, handlerConfig := range davHandlerConfigs {
		log.Debug().Interface("handlerConfig", handlerConfig).Msg("dav")
	}

	handlerConfigs = append(handlerConfigs, davHandlerConfigs...)
	return handlerConfigs
}

func Execute() {
	type Config struct {
		Address         string `koanf:"address"`
		Port            string `koanf:"port"`
		Dav             string `koanf:"dav" usage:"dav config, format: prefix,pathDir,username,password,readonly;..."`
		DavListIsSecret bool   `koanf:"secret_dav_list" usage:"if true, hide the dav list"`
		PreRequestHook string `koanf:"pre_request_hook" usage:"path to the pre request hook"`
	}
	cfg := &Config{
		Address:         "0.0.0.0",
		Port:            "80",
		Dav:             DEFAULT_DAV_CONFIG,
		DavListIsSecret: false,
		PreRequestHook:  "",
	}
	result := goutils.LoadConfig(cfg)
	log.Info().Interface("config", cfg).Msg("Config")

	handlerConfigs := getHandlerConfigs(cfg.Dav, result.K.Slices("davs"))

	if isDefaultHandlerConfigs(handlerConfigs) {
		createDefaultDirs()
		log.Info().Msg("Default dav config is used, created default directories")
	}

	server, err := server.NewWebDAVServer(cfg.Address+":"+cfg.Port, handlerConfigs, cfg.DavListIsSecret,
		cfg.PreRequestHook)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create server")
	}
	err = server.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to run server")
	}
}
