package config

import (
	"io/ioutil"
	"log"
	"os"
)

var Log *Logging

type Logging struct {
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
}

func init() {
	Log = NewLogging(true)
}

func NewLogging(enabled bool) *Logging {

	var infoLog *log.Logger
	var errorLog *log.Logger

	if enabled {
		infoLog = log.New(openLogFile("info.log"), "INFO\t", log.Ldate|log.Ltime)
		errorLog = log.New(openLogFile("error.log"), "ERROR\t", log.Ldate|log.Ltime)
	} else {
		infoLog = log.New(ioutil.Discard, "", 0)
		errorLog = log.New(ioutil.Discard, "", 0)
	}

	app := &Logging{
		InfoLog:    infoLog,
		ErrorLog: errorLog,
	}

	return app
}

func openLogFile(path string) *os.File {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	return f
}
