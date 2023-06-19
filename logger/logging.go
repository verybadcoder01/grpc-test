package logger

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
)

func SetupLogging(path string, maxSize int, maxAge int, maxBackups int) {
	_, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    maxSize, // megabytes
		MaxBackups: maxBackups,
		MaxAge:     maxAge, //days
		Compress:   true,   // disabled by default
	})
}
