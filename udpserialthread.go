package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/goburrow/serial"
)

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
	serialPortOptions := serial.Config{
		Address:  ttyName,
		BaudRate: portConfig.BaudRate,
		DataBits: portConfig.DataBits,
		StopBits: portConfig.StopBits,
		Parity:   "N",
		Timeout:  4 * time.Millisecond,
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
	serialPort, err := serial.Open(&serialPortOptions)
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
	/*
		var serial2udpQueue = lane.NewQueue()
		defer func() {
			for {
				if serial2udpQueue.Empty() == false {
					serial2udpQueue.Dequeue()
				} else {
					break
				}
			}
		}()
	*/

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
		separator := parsePacketSeparator(portConfig.PacketSeparator)
		var readLength int
		tempSerialBuffer := make([]byte, 1)
		for {
			readLength, err = serialPort.Read(tempSerialBuffer)
			if err != serial.ErrTimeout && err != nil {
				logger(name, LogWarning, err)
			} else {
				packetOut := false
				if readLength > 0 {
					serialBuffer[serialBufferContentSize] = tempSerialBuffer[0]
					serialBufferContentSize++
					if serialBufferContentSize >= len(serialBuffer) {
						packetOut = true
					}
					if (portConfig.PacketSeparator != "") && (string(tempSerialBuffer[0]) == separator) {
						if PrintDebug {
							fmt.Println("Hit separator")
						}
						packetOut = true
					}
				} else {
					if serialBufferContentSize > 0 {
						packetOut = true
					}
					time.Sleep(time.Nanosecond * 100)
				}
				if packetOut == true {
					//serial2udpQueue.Enqueue(serialBuffer[:serialBufferContentSize])
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
