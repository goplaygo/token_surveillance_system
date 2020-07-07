/*
@Time : 2020/7/7 1:30 下午
@Author : L
@File : logger.go
@Software: GoLand
*/
package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"token_surveillance_system/extend/conf"
)

func Setup(){
	switch strings.ToLower(conf.LoggerConf.Level){
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	if conf.LoggerConf.Pretty {
		 log.Logger = log.Output(zerolog.ConsoleWriter{
		 	Out: os.Stderr,
		 	NoColor: !conf.LoggerConf.Color,
		 })
	}
}
