package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "GoWebDAV/internal/common"
	"GoWebDAV/internal/server"

	"github.com/spf13/cobra"
)

const DEFAULT_DAV_CONFIG = "/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false"

var (
	// Used for flags.
	dav  string
	addr string
	port string

	rootCmd = &cobra.Command{
		Use: "GoWebDAV",
		Run: func(cmd *cobra.Command, args []string) {
			if dav == DEFAULT_DAV_CONFIG {
				createDefaultDirs()
			}

			handlerConfigs, err := parseDavToHandlerConfigs(dav)
			if err != nil {
				panic(err)
			}
			server, err := server.NewWebDAVServer(addr+":"+port, handlerConfigs)
			if err != nil {
				panic(err)
			}
			server.Run()
		},
	}
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
		arr := strings.Split(davConfig, ",")
		if len(arr) != 5 {
			err = fmt.Errorf("invalid dav config: %s", davConfig)
			return
		}
		prefix := arr[0]
		pathDir := arr[1]
		username := arr[2]
		password := arr[3]
		readonly, err := strconv.ParseBool(arr[4])
		if err != nil {
			readonly = false
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

// Execute executes the root command.
func Execute() error {
	// fmt.Println(rootCmd.UsageTemplate())
	rootCmd.SetUsageTemplate("Flags:\n{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}")

	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println(rootCmd.UsageString())
	})
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dav, "dav", DEFAULT_DAV_CONFIG, "")
	rootCmd.PersistentFlags().StringVarP(&addr, "address", "a", "0.0.0.0", "address to listen")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "80", "port to listen")
}
