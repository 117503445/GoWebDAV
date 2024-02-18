package common

import (
	// "fmt"
	// "strings"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	setLogger()
}

func setLogger() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}
	log.Logger = log.Output(output)
}
