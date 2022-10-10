/*
  udpserial - lightweight bridge for serial ports over UDP packets

  Copyright (C) 2022  Giacomo De Lazzari

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

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
