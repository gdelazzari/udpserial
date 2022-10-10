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
	"encoding/json"
	"errors"
	"io/ioutil"
)

// PortDefinition : structure that contains information about the mapping of a port name to a tty
type PortDefinition struct {
	PortName string `json:"name"`
	TTY      string `json:"tty"`
}

// Definitions : structure that represents the definitions file data
type Definitions struct {
	PortDefinitions []PortDefinition `json:"ports"`
	BaudRates       []int            `json:"baudrates"`
}

func readDefinitions(filename string) Definitions {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		logger("definitions", LogFatal, err)
	}

	var definitions Definitions
	json.Unmarshal(raw, &definitions)
	return definitions
}

func getPortTTY(definitions Definitions, portname string) (string, error) {
	for _, portDefinition := range definitions.PortDefinitions {
		if portDefinition.PortName == portname {
			return portDefinition.TTY, nil
		}
	}

	return "", errors.New("no such port name " + portname + " in definitions file")
}
