package main

import (
	"sync"
	"time"
)

const definitionsFilename = "definitions.json"
const configFilename = "config.json"

var definitions Definitions
var config Config

func main() {
	initLogger()
	defer closeLogger()

	definitions = readDefinitions(definitionsFilename)
	logger("main", LogInfo, "Loaded definitions")

	config = readConfig(configFilename)
	logger("main", LogInfo, "Loaded configuration")

	// TODO check if ports n > 0

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go startAndSuperviseThreads(&waitGroup)

	waitGroup.Add(1)
	go serveWebPanel(&waitGroup)

	time.Sleep(time.Second * 3)

	go restartAllThreads()

	time.Sleep(time.Millisecond * 250)

	restartAllThreads()

	waitGroup.Wait()
}
