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
	"strconv"
	"sync"
	"time"
)

var stopChannel chan string
var killChannels map[string](chan bool)

var doNotRestart bool
var diedCount int

var restarting bool

func rebuildStatistics() {
	statistics.PortsMutex.Lock()

	statistics.Ports = make(map[string]*PortStatistics)

	for _, portConfig := range config.Ports {
		statistics.Ports[portConfig.Name] = &PortStatistics{}
	}

	statistics.PortsMutex.Unlock()
}

func startAndSuperviseThreads(wg *sync.WaitGroup) {
	defer wg.Done()

	logger("supervisor", LogInfo, "Starting threads")

	restarting = false
	doNotRestart = false
	diedCount = 0

	stopChannel = make(chan string)
	killChannels = make(map[string](chan bool))

	rebuildStatistics()

	for _, portConfig := range config.Ports {
		name := "UDPSerialThread_" + portConfig.Name
		killChannels[portConfig.Name] = make(chan bool)
		go UDPSerialThread(name, portConfig, stopChannel, killChannels[portConfig.Name], statistics.Ports[portConfig.Name])
	}

	// Start supervising
	for {
		var diedPortName = <-stopChannel

		diedCount++

		if doNotRestart == false {
			portConfig, err := getPortConfig(config, diedPortName)
			if err != nil {
				logger("supervisor", LogError, err)
			}

			// Relaunch thread
			name := "UDPSerialThread_" + portConfig.Name
			killChannels[portConfig.Name] = make(chan bool)
			go UDPSerialThread(name, portConfig, stopChannel, killChannels[portConfig.Name], statistics.Ports[portConfig.Name])

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

	logger("supervisor", LogInfo, "Waiting for "+strconv.Itoa(mustDie)+" threads to stop")

	for portName := range killChannels {
		select {
		case killChannels[portName] <- true:
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
	rebuildStatistics()

	for _, portConfig := range config.Ports {
		name := "UDPSerialThread_" + portConfig.Name
		killChannels[portConfig.Name] = make(chan bool)
		go UDPSerialThread(name, portConfig, stopChannel, killChannels[portConfig.Name], statistics.Ports[portConfig.Name])
	}

	restarting = false
}
