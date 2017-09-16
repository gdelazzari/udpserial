package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// PortConfig : structure holding a port configuration parameters
type PortConfig struct {
	Name            string `json:"name"`
	BaudRate        uint   `json:"baudrate"`
	DataBits        uint   `json:"databits"`
	StopBits        uint   `json:"stopbits"`
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
