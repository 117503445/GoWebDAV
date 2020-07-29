package model

import (
	"golang.org/x/net/webdav"
)

type WebDAVConfig struct {
	Prefix   string
	PathDir  string
	Username string
	Password string
	Handler  *webdav.Handler
}

func (WebDAVConfig *WebDAVConfig) Init(prefix string, pathDir string, username string, password string) {
	WebDAVConfig.Prefix = prefix
	WebDAVConfig.PathDir = pathDir
	WebDAVConfig.Username = username
	WebDAVConfig.Password = password

	WebDAVConfig.Handler = &webdav.Handler{
		FileSystem: webdav.Dir(pathDir),
		LockSystem: webdav.NewMemLS(),
		Prefix:     prefix,
	}
}
func WebDAVConfigFindOneByPrefix(WebDAVConfigs []*WebDAVConfig, prefix string) *WebDAVConfig {
	for _, WebDAVConfig := range WebDAVConfigs {
		if WebDAVConfig.Prefix == prefix {
			return WebDAVConfig
		}
	}
	return nil
}
