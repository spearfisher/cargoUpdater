package utils

import (
	"log"
	"os"
)

// Logger is a application file logger
var Logger *log.Logger

func init() {
	logpath := os.Getenv("DELLA_LOGPATH")
	if logpath == "" {
		logpath = "./della_log.txt"
	}

	file, err := os.OpenFile(logpath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	Logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	Logger.Println("Application started")
}
