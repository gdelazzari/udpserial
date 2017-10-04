package main

import (
	"sync"
)

import _ "net/http/pprof"

const definitionsFilename = "definitions.json"
const configFilename = "config.json"

var definitions Definitions
var config Config

var statistics Statistics

func main() {
	initLogger()
	defer closeLogger()

	definitions = readDefinitions(definitionsFilename)
	logger("main", LogInfo, "Loaded definitions")

	config = readConfig(configFilename)
	logger("main", LogInfo, "Loaded configuration")

	if len(config.Ports) <= 0 {
		logger("main", LogWarning, "No ports configured")
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go startAndSuperviseThreads(&waitGroup)

	waitGroup.Add(1)
	go serveWebPanel(&waitGroup)

	waitGroup.Add(1)
	go statisticsThread(&waitGroup, &statistics)

	waitGroup.Wait()
}
