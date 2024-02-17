package model

import (
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/webdav"
)

type WebDAVConfig struct {
	Prefix   string
	PathDir  string
	Username string
	Password string
	ReadOnly bool
	Handler  *webdav.Handler
}

func (WebDAVConfig *WebDAVConfig) Init(prefix string, pathDir string, username string, password string, readonly bool) {
	WebDAVConfig.Prefix = prefix
	WebDAVConfig.PathDir = pathDir
	WebDAVConfig.Username = username
	WebDAVConfig.Password = password
	WebDAVConfig.ReadOnly = readonly

	WebDAVConfig.Handler = &webdav.Handler{
		FileSystem: webdav.Dir(pathDir),
		LockSystem: webdav.NewMemLS(),
		Prefix:     prefix,
	}
}

func InitByConfigStr(str string) WebDAVConfig {
	davConfigArray := strings.Split(str, ",")
	prefix := davConfigArray[0]
	pathDir := davConfigArray[1]
	username := davConfigArray[2]
	password := davConfigArray[3]

	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
		println("Dir not exists: ", pathDir)
		os.Exit(1)
	}

	readonly, err := strconv.ParseBool(davConfigArray[4])
	if err != nil {
		readonly = false
	}

	var WebDAVConfig WebDAVConfig
	WebDAVConfig.Init(prefix, pathDir, username, password, readonly)
	return WebDAVConfig
}

func ParseURL(WebDAVConfigs []*WebDAVConfig, url string) (*WebDAVConfig, string) {
	for _, WebDAVConfig := range WebDAVConfigs {
		davPath, found := strings.CutPrefix(url, WebDAVConfig.Prefix)
		if found {
			return WebDAVConfig, davPath
		}
	}
	return nil, ""
}
