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

// PortConfig : structure holding a port configuration parameters
type PortConfig struct {
	Name            string `json:"name"`
	BaudRate        int    `json:"baudrate"`
	DataBits        int    `json:"databits"`
	StopBits        int    `json:"stopbits"`
	PacketSeparator string `json:"packetSeparator"`
	UDPInputIP      string `json:"udpInputIP"`
	UDPInputPort    int    `json:"udpInputPort"`
	UDPOutputIP     string `json:"udpOutputIP"`
	UDPOutputPort   int    `json:"udpOutputPort"`
}

// Config : structure holding the service configuration parameters
type Config struct {
	Ports []PortConfig `json:"ports"`
}

func (config *Config) toJSON() string {
	bytes, err := config.toJSONBytes()
	if err != nil {
		return ""
	}

	return string(bytes)
}

func (config *Config) toJSONBytes() ([]byte, error) {
	return json.Marshal(config)
}

func (config *Config) saveToFile(filename string) {
	bytes, err := config.toJSONBytes()
	if err != nil {
		logger("config", LogError, err)
	} else {
		err := ioutil.WriteFile(filename, bytes, 0644)
		if err != nil {
			logger("config", LogError, err)
		}
	}
}

func readConfig(filename string) Config {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		logger("config", LogWarning, err)
		logger("config", LogInfo, "creating empty configuration file")
		var config = Config{}
		config.saveToFile(filename)
		return Config{}
	}

	var config Config
	json.Unmarshal(raw, &config)
	return config
}

func getPortConfig(config Config, portname string) (PortConfig, error) {
	for _, portConfig := range config.Ports {
		if portConfig.Name == portname {
			return portConfig, nil
		}
	}
	return PortConfig{}, errors.New("no such port " + portname + " in configuration file")
}
