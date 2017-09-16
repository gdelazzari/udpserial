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
