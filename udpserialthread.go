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
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

// PrintDebug : print debug information while running
const PrintDebug = false

func findPacketSeparator(buffer []byte, separator string) int {
	for i, b := range buffer {
		if b == separator[0] {
			return i
		}
	}
	return -1
}

func parsePacketSeparator(ps string) string {
	result := ps
	result = strings.Replace(result, "\\n", "\n", -1)
	result = strings.Replace(result, "\\r", "\r", -1)
	result = strings.Replace(result, "\\0", "\000", -1)
	result = strings.Replace(result, "\\t", "\t", -1)
	return result
}

// UDPSerialThread : start a loop for a specified port
func UDPSerialThread(name string, portConfig PortConfig, stopChannel chan string, killChannel chan bool, stats *PortStatistics) {
	defer func() { stopChannel <- portConfig.Name }()

	logger(name, LogInfo, "Starting thread")

	var running = true

	// Get serial port TTY
	ttyName, err := getPortTTY(definitions, portConfig.Name)
	if err != nil {
		logger(name, LogError, err)
		stats.Errors++
		return
	}

	// Serial port configuration
	serialPortOptions := serial.OpenOptions{
		PortName:              ttyName,
		BaudRate:              uint(portConfig.BaudRate),
		DataBits:              uint(portConfig.DataBits),
		StopBits:              uint(portConfig.StopBits),
		MinimumReadSize:       0,
		InterCharacterTimeout: 100,
	}

	// UDP Input Address
	udpInputAddress, err := net.ResolveUDPAddr("udp", portConfig.UDPInputIP+":"+strconv.Itoa(portConfig.UDPInputPort))
	if err != nil {
		logger(name, LogError, err)
		stats.Errors++
		return
	}

	// UDP Output Address
	udpOutputAddress := portConfig.UDPOutputIP + ":" + strconv.Itoa(portConfig.UDPOutputPort)

	// Open serial port
	serialPort, err := serial.Open(serialPortOptions)
	if err != nil {
		logger(name, LogError, err)
		stats.Errors++
		return
	}
	logger(name, LogInfo, "Opened "+ttyName)
	defer serialPort.Close()

	// Open UDP input connection
	udpInputConnection, err := net.ListenUDP("udp", udpInputAddress)
	if err != nil {
		logger(name, LogError, err)
		stats.Errors++
		return
	}
	logger(name, LogInfo, "Listening on "+udpInputAddress.String())
	defer udpInputConnection.Close()

	// Open UDP output connection
	udpOutputConnection, err := net.Dial("udp", udpOutputAddress)
	if err != nil {
		logger(name, LogError, err)
		stats.Errors++
		return
	}
	logger(name, LogInfo, "Sending to "+udpOutputAddress)
	defer udpOutputConnection.Close()

	var udpBuffer = make([]byte, 5100)
	var serialBuffer = make([]byte, 4096)
	var serialBufferContentSize = 0

	var serial2udpChannel = make(chan []byte, 64)

	var serialChannel = make(chan byte, 1024)

	var internalWaitGroup sync.WaitGroup

	internalWaitGroup.Add(1)
	go func() {
		defer internalWaitGroup.Done()
		var readLength int
		//var readAddr *net.UDPAddr
		for {
			readLength, _, err = udpInputConnection.ReadFromUDP(udpBuffer)
			if err != nil {
				logger(name, LogWarning, err)
			} else {
				stats.UDP2SerialCounter += readLength
				if PrintDebug {
					fmt.Println("UDP in ", udpBuffer[0:readLength])
				}
				go func() {
					if PrintDebug {
						fmt.Println("writing to serial port")
					}
					_, err = serialPort.Write(udpBuffer[:readLength])
					if err != nil {
						// TODO do not repeat error for every packet
						logger(name, LogWarning, err)
					} else {
						if PrintDebug {
							fmt.Println("wrote to serial port")
						}
					}
				}()
			}
			if running == false {
				logger(name, LogInfo, "udp2serial subthread stopped")
				break
			}
		}
	}()

	internalWaitGroup.Add(1)
	go func() {
		defer internalWaitGroup.Done()

		tempSerialBuffer := make([]byte, 256)

		for {
			readLength, err := serialPort.Read(tempSerialBuffer)

			if err == nil && readLength > 0 {
				for _, b := range tempSerialBuffer[:readLength] {
					serialChannel <- b
				}
			}

			if running == false {
				logger(name, LogInfo, "serialReader subthread stopped")
				break
			}
		}
	}()

	internalWaitGroup.Add(1)
	go func() {
		defer internalWaitGroup.Done()
		separator := parsePacketSeparator(portConfig.PacketSeparator)
		for {
			packetOut := false

			// Timeout read from serialChannel (for an incoming serial byte)
			select {
			case incoming := <-serialChannel:
				serialBuffer[serialBufferContentSize] = incoming
				serialBufferContentSize++
				if serialBufferContentSize >= len(serialBuffer) {
					packetOut = true
				}
				if (portConfig.PacketSeparator != "") && (string(incoming) == separator) {
					if PrintDebug {
						fmt.Println("Hit separator")
					}
					packetOut = true
				}
			case <-time.After(5 * time.Millisecond):
				// If nothing arrives within the specified time, flush out the buffer to UDP
				packetOut = true
			}

			if packetOut == true {
				if serialBufferContentSize > 0 {
					serial2udpChannel <- serialBuffer[:serialBufferContentSize]
					stats.Serial2UDPCounter += serialBufferContentSize
					sizebefore := serialBufferContentSize
					serialBufferContentSize = 0
					if PrintDebug {
						fmt.Println("Ready out: ", serialBuffer[:sizebefore])
					}
				}
			}

			if running == false {
				logger(name, LogInfo, "serial2udpqueue subthread stopped")
				break
			}
		}
	}()

	internalWaitGroup.Add(1)
	timeoutChannel := make(chan bool, 1)
	go func() {
		defer internalWaitGroup.Done()
		for {
			time.Sleep(100 * time.Millisecond)
			timeoutChannel <- true
			if running == false {
				break
			}
		}
	}()

	internalWaitGroup.Add(1)
	go func() {
		defer internalWaitGroup.Done()
		for {
			select {
			case toSend := <-serial2udpChannel:
				if PrintDebug {
					fmt.Println("UDP out ", toSend)
				}
				_, err = udpOutputConnection.Write(toSend)
				if err != nil {
					// TODO do not repeat error for every packet
					// logger(name, LogWarning, err)
					if PrintDebug {
						fmt.Printf("UDP refused for packet %q\n", toSend)
					}
					stats.LostPackets++
				} else {
					if PrintDebug {
						fmt.Printf("UDP sent for packet %q\n", toSend)
					}
				}
			case <-timeoutChannel:
				// just move on
			}
			/*
				if serial2udpQueue.Empty() == false {
					toSend, assertResult = serial2udpQueue.Dequeue().([]byte)
					if assertResult == false {
						logger(name, LogWarning, "serial2udpQueue dequeued a non-[]byte node")
					} else {
						if PrintDebug {
							fmt.Println("UDP out ", toSend)
						}
						_, err = udpOutputConnection.Write(toSend)
						if err != nil {
							// TODO do not repeat error for every packet
							// logger(name, LogWarning, err)
							if PrintDebug {
								fmt.Printf("UDP refused for packet %q\n", toSend)
							}
							stats.LostPackets++
						} else {
							if PrintDebug {
								fmt.Printf("UDP sent for packet %q\n", toSend)
							}
						}
					}
				}
			*/
			if running == false {
				logger(name, LogInfo, "udpqueue2udp subthread stopped")
				break
			}
			//time.Sleep(time.Nanosecond * 1000)
		}
	}()

	internalWaitGroup.Add(1)
	go func() {
		defer internalWaitGroup.Done()
		<-killChannel
		logger(name, LogInfo, "Thread received kill signal")
		serialPort.Close()
		udpInputConnection.Close()
		udpOutputConnection.Close()
		running = false
	}()

	internalWaitGroup.Wait()

	logger(name, LogWarning, "Thread reached end")
}
