package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var logFile *os.File

// Log level constants
const (
	LogInfo    = 4
	LogWarning = 3
	LogError   = 2
	LogFatal   = 1
	LogPanic   = 0
)

func initLogger() {
	var err error
	_ = os.Remove("udpserial.log.bak")
	_ = os.Rename("udpserial.log", "udpserial.log.bak")
	logFile, err = os.OpenFile("udpserial.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	multiwriter := io.MultiWriter(os.Stdout, logFile)

	log.SetOutput(multiwriter)
}

func logger(source string, level int, content interface{}) {
	switch level {
	case LogInfo:
		log.Println(("[" + source + "]"), "[INFO]", content)
	case LogWarning:
		log.Println(("[" + source + "]"), "[WARNING]", content)
	case LogError:
		log.Println(("[" + source + "]"), "[ERROR]", content)
	case LogFatal:
		log.Fatalln(("[" + source + "]"), "[FATAL]", content)
	case LogPanic:
		log.Panicln(("[" + source + "]"), "[PANIC]", content)
	}
}

func getLogString() string {
	bytes, err := ioutil.ReadFile("udpserial.log")
	if err != nil {
		logger("logger", LogError, err)
		return ""
	}
	return string(bytes)
}

func closeLogger() {
	logFile.Close()
}
