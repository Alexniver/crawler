package crawler

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/rifflock/lfshook"
	"os"
)

var _logger *log.Logger

var DefaultLogLevelFilePathMap = lfshook.PathMap{
	log.InfoLevel:  "_info.log",
	log.DebugLevel: "_debug.log",
	log.WarnLevel:  "_warn.log",
	log.ErrorLevel: "_error.log",
	log.FatalLevel: "_fatal.log",
	log.PanicLevel: "_panic.log",
}

func GetLogger(logLevelFilePathMap map[log.Level]string) *log.Logger {
	if _logger != nil {
		return _logger
	}

	//create if not exist
	for _, path := range logLevelFilePathMap {
		if _, err := os.Stat(path); os.IsNotExist(err) {

			file, err := os.Create(path)
			defer file.Close()
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	_logger = log.New()
	//_logger.Formatter = new(log.JSONFormatter)
	_logger.Formatter = new(log.TextFormatter)
	_logger.Hooks.Add(lfshook.NewHook(logLevelFilePathMap))

	return _logger
}

func GetDefaultLogger() *log.Logger {
	return GetLogger(DefaultLogLevelFilePathMap)
}
