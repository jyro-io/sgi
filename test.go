package sgi

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.ErrorLevel)
}

func main() {
	contextLogger := log.WithFields(log.Fields{"function": "main"})
	s := NewClient(
		5,
		"localhost",
		"http",
		"test",
		"test")
	r, err := GetDefinition(s, "archimedes", "datasource", "test")
	if err == nil {
		contextLogger.Fatal(err)
		os.Exit(1)
	} else {
		if r.status == true {
			contextLogger.Info(r.response)
		} else {
			contextLogger.Fatal(r.response)
		}
	}
}