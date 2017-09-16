package main

import (
	"sync"
	"time"
)

var stopChannel chan string
var killChannels map[string](chan bool)

var doNotRestart bool
var diedCount int

var restarting bool

func startAndSuperviseThreads(wg *sync.WaitGroup) {
	defer wg.Done()

	logger("supervisor", LogInfo, "Starting threads")

	restarting = false
	doNotRestart = false
	diedCount = 0

	stopChannel = make(chan string)
	killChannels = make(map[string](chan bool))

	for _, portConfig := range config.Ports {
		name := "UDPSerialThread_" + portConfig.Name
		killChannels[portConfig.Name] = make(chan bool)
		go UDPSerialThread(name, portConfig, stopChannel, killChannels[portConfig.Name])
	}

	// Start supervising
	for {
		var diedPortName = <-stopChannel

		diedCount++

		if doNotRestart == false {
			portConfig, err := getPortConfig(config, diedPortName)
			if err != nil {
				logger("supervisor", LogFatal, err)
			}

			// Relaunch thread
			name := "UDPSerialThread_" + portConfig.Name
			killChannels[portConfig.Name] = make(chan bool)
			go UDPSerialThread(name, portConfig, stopChannel, killChannels[portConfig.Name])

			// Wait a bit
			time.Sleep(time.Second * 1)
		}
	}
}

func stopAllThreads() {
	if restarting == true {
		return
	}

	restarting = true

	doNotRestart = true
	diedCount = 0

	mustDie := len(killChannels)

	for _, portConfig := range config.Ports {
		select {
		case killChannels[portConfig.Name] <- true:
		case <-time.After(1*time.Second + time.Millisecond*250):
		}
	}

	for {
		if diedCount >= mustDie {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}

	doNotRestart = false

	logger("supervisor", LogInfo, "All threads stopped")
}

func restartAllThreads() {
	logger("supervisor", LogInfo, "Thread restart requested")
	if restarting == true {
		logger("supervisor", LogInfo, "Already restarting threads")
		return
	}

	stopAllThreads()

	restarting = true

	logger("supervisor", LogInfo, "Restarting threads")

	killChannels = make(map[string](chan bool))

	for _, portConfig := range config.Ports {
		name := "UDPSerialThread_" + portConfig.Name
		killChannels[portConfig.Name] = make(chan bool)
		go UDPSerialThread(name, portConfig, stopChannel, killChannels[portConfig.Name])
	}

	restarting = false
}
