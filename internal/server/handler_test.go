package server

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const UnExistDir = "/114514"

func TestMain(m *testing.M) {
	if _, err := os.Stat(UnExistDir); err == nil {
		panic(UnExistDir + " should not exist")
	}

	code := m.Run()
	os.Exit(code)
}

func TestHandlerConfig(t *testing.T) {
	assert := assert.New(t)

	wd, err := os.Getwd()
	assert.Nil(err)

	cases := []struct {
		cfg   *HandlerConfig
		valid bool
	}{
		{
			&HandlerConfig{
				Prefix:  "/data",
				PathDir: wd,
			},
			true,
		}, {
			&HandlerConfig{
				Prefix:  "/",
				PathDir: wd,
			},
			true,
		}, {
			&HandlerConfig{
				Prefix:  "data",
				PathDir: wd,
			},
			false,
		}, {
			&HandlerConfig{
				Prefix:  "",
				PathDir: wd,
			},
			false,
		}, {
			&HandlerConfig{
				Prefix:  "/data/",
				PathDir: wd,
			},
			false,
		}, {
			&HandlerConfig{
				Prefix:  "/114?514",
				PathDir: wd,
			},
			false,
		}, {
			&HandlerConfig{
				Prefix:  "/data",
				PathDir: "/114514",
			},
			false,
		}, {
			&HandlerConfig{
				Prefix:   "/data",
				PathDir:  wd,
				Username: "user",
			},
			false,
		},
	}

	for _, c := range cases {
		if c.valid {
			assert.Nil(checkHandlerConfig(c.cfg, false), "cfg: %+v should be valid", c.cfg)
		} else {
			assert.NotNil(checkHandlerConfig(c.cfg, false), "cfg: %+v should be invalid", c.cfg)
		}
	}
}

func TestHandlerConfigs(t *testing.T) {
	assert := assert.New(t)

	wd, err := os.Getwd()
	assert.Nil(err)

	cases := []struct {
		cfgs  []*HandlerConfig
		valid bool
	}{
		{
			[]*HandlerConfig{
				{
					Prefix:  "/data1",
					PathDir: wd,
				},
				{
					Prefix:  "/data2",
					PathDir: wd,
				},
			},
			true,
		}, {
			[]*HandlerConfig{
				{
					Prefix:  "/data1",
					PathDir: wd,
				},
				{
					Prefix:  "/data1",
					PathDir: wd,
				},
			},
			false,
		},
	}

	for _, c := range cases {
		if c.valid {
			assert.Nil(checkHandlerConfigs(c.cfgs), "cfgs: %+v should be valid", c.cfgs)
		} else {
			assert.NotNil(checkHandlerConfigs(c.cfgs), "cfgs: %+v should be invalid", c.cfgs)
		}
	}
}
