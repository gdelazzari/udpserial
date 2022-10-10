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
	"time"
)

// PortStatistics : represent the statistics for one port
type PortStatistics struct {
	UDP2SerialRate    int
	Serial2UDPRate    int
	LostPackets       int
	Errors            int
	UDP2SerialCounter int
	Serial2UDPCounter int
}

// PublicPortStatistics : represent the public information about a port statistics
type PublicPortStatistics struct {
	UDP2SerialRate int `json:"udp2serialRate"`
	Serial2UDPRate int `json:"serial2udpRate"`
	LostPackets    int `json:"lostPackets"`
	Errors         int `json:"errors"`
}

// Statistics : represent statistics for all ports
type Statistics struct {
	Ports      map[string]*PortStatistics
	PortsMutex sync.Mutex
}

// PublicStatistics : represent the public statistics information
type PublicStatistics struct {
	Ports map[string]PublicPortStatistics `json:"ports"`
}

func getPublicStatistics(stats *Statistics) PublicStatistics {
	result := PublicStatistics{}
	result.Ports = make(map[string]PublicPortStatistics)

	stats.PortsMutex.Lock()

	for portName := range stats.Ports {
		result.Ports[portName] = PublicPortStatistics{
			stats.Ports[portName].UDP2SerialRate,
			stats.Ports[portName].Serial2UDPRate,
			stats.Ports[portName].LostPackets,
			stats.Ports[portName].Errors,
		}
	}

	stats.PortsMutex.Unlock()

	return result
}

func statisticsThread(wg *sync.WaitGroup, stats *Statistics) {
	defer wg.Done()

	for {
		time.Sleep(time.Second * 1)

		stats.PortsMutex.Lock()

		for portName := range stats.Ports {
			stats.Ports[portName].Serial2UDPRate = stats.Ports[portName].Serial2UDPCounter
			stats.Ports[portName].UDP2SerialRate = stats.Ports[portName].UDP2SerialCounter
			stats.Ports[portName].Serial2UDPCounter = 0
			stats.Ports[portName].UDP2SerialCounter = 0
		}

		stats.PortsMutex.Unlock()
	}
}
