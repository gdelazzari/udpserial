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
	"sync"
)

import _ "net/http/pprof"

const definitionsFilename = "definitions.json"
const configFilename = "config.json"

var definitions Definitions
var config Config

var statistics Statistics

func main() {
	println(`udpserial  Copyright (C) 2022  Giacomo De Lazzari

This program comes with ABSOLUTELY NO WARRANTY.
This is free software, and you are welcome to redistribute it
under certain conditions; see the provided LICENSE for details.
`)
    
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
