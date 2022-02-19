package logger

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Init() error {
	if viper.GetString("log.format") == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetReportCaller(true)
		logrus.SetFormatter(&logrus.TextFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				modDotPos := strings.LastIndex(filename, ".")
				funcName := path.Base(f.Function)
				funcDotPos := strings.LastIndex(funcName, ".")
				return fmt.Sprintf("[%s/%s]", filename[:modDotPos], funcName[funcDotPos+1:]), ""
			},
		})
	}

	level := logrus.TraceLevel
	switch viper.GetString("log.level") {
	case "error":
		level = logrus.ErrorLevel
	case "warn":
		level = logrus.WarnLevel
	case "info":
		level = logrus.InfoLevel
	case "debug":
		level = logrus.DebugLevel
	case "trace":
		level = logrus.TraceLevel
	}
	logrus.SetLevel(level)

	return nil
}
