package model

import (
	"golang.org/x/net/webdav"
	"strconv"
	"strings"
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
func (WebDAVConfig *WebDAVConfig) InitByConfigStr(str string) {
	davConfigArray := strings.Split(str, ",")
	prefix := davConfigArray[0]
	pathDir := davConfigArray[1]
	username := davConfigArray[2]
	password := davConfigArray[3]

	readonly, err := strconv.ParseBool(davConfigArray[4])
	if err != nil {
		readonly = false
	}

	WebDAVConfig.Init(prefix, pathDir, username, password, readonly)
}
func WebDAVConfigFindOneByPrefix(WebDAVConfigs []*WebDAVConfig, prefix string) *WebDAVConfig {
	for _, WebDAVConfig := range WebDAVConfigs {
		if WebDAVConfig.Prefix == prefix {
			return WebDAVConfig
		}
	}
	return nil
}
