package crawler

import (
	log "github.com/Sirupsen/logrus"
	//"github.com/rifflock/lfshook"
	"testing"
)

func TestNewLogger(*testing.T) {
	var l = GetDefaultLogger()
	l.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")
	l.Info("haha")
	l.Error("test error")
}
