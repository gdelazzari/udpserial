package main

import (
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/goburrow/serial"
	"github.com/oleiade/lane"
)

func findPacketSeparator(buffer []byte, separator string) int {
	for i, b := range buffer {
		if b == separator[0] {
			return i
		}
	}
	return -1
}

func transferBuffers(src []byte, srcStart int, dest *[]byte, destStart int, len int) {
	for i := 0; i < len; i++ {
		(*dest)[destStart+i] = src[srcStart+i]
	}
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
	/*
		if err != nil {
			logger(name, LogError, err)
			stats.Errors++
			return
		}
	*/

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
				//fmt.Println("[", name, "] UDP: ", string(udpBuffer[0:readLength]), " from ", readAddr)
				stats.UDP2SerialCounter += readLength
				go func() {
					_, err = serialPort.Write(udpBuffer[:readLength])
					if err != nil {
						// TODO do not repeat error for every packet
						logger(name, LogWarning, err)
					} else {
						//fmt.Println(name, "wrote to serial port")
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
		var readLength int
		tempSerialBuffer := make([]byte, 1024)
		for {
			readLength, err = serialPort.Read(tempSerialBuffer)
			if err != serial.ErrTimeout && err != nil {
				logger(name, LogWarning, err)
			} else {
				/*if readLength > 0 {
					stats.Serial2UDPCounter += readLength
					fmt.Printf("Serial: %q\n", serialBuffer[:readLength])
					serial2udpQueue.Enqueue(serialBuffer[:readLength])
				}*/
				if readLength > 0 {
					/*
						if len(portConfig.PacketSeparator) > 0 {
							packetSeparatorPos := findPacketSeparator(serialBuffer, portConfig.PacketSeparator)
							if packetSeparatorPos >= 0 {
								// TODO check main buffer overflow
								transferBuffers(tempSerialBuffer, 0, &serialBuffer, serialBufferContentSize, packetSeparatorPos+1)
								serialBufferContentSize += packetSeparatorPos + 1

								serial2udpQueue.Enqueue(serialBuffer[:serialBufferContentSize])
								serialBufferContentSize = 0

								len := readLength - packetSeparatorPos - 1
								transferBuffers(tempSerialBuffer, packetSeparatorPos+1, &serialBuffer, serialBufferContentSize, len)
								serialBufferContentSize += len
							} else {
								transferBuffers(tempSerialBuffer, 0, &serialBuffer, serialBufferContentSize, readLength)
								serialBufferContentSize += readLength
							}
						} else {
					*/
					//fmt.Printf("Serial: %q\n", tempSerialBuffer[:readLength])
					if (serialBufferContentSize + readLength) >= len(serialBuffer) {
						serial2udpQueue.Enqueue(serialBuffer[:serialBufferContentSize])
						//fmt.Printf("Ready out %q\n", serialBuffer[:serialBufferContentSize])
						stats.Serial2UDPCounter += serialBufferContentSize
						serialBufferContentSize = 0
					}
					transferBuffers(tempSerialBuffer, 0, &serialBuffer, serialBufferContentSize, readLength)
					serialBufferContentSize += readLength
					//}
				} else {
					if serialBufferContentSize > 0 {
						serial2udpQueue.Enqueue(serialBuffer[:serialBufferContentSize])
						//fmt.Printf("Ready out %q\n", serialBuffer[:serialBufferContentSize])
						stats.Serial2UDPCounter += serialBufferContentSize
						serialBufferContentSize = 0
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
	go func() {
		defer internalWaitGroup.Done()
		for {
			if serial2udpQueue.Empty() == false {
				var toSend, assertResult = serial2udpQueue.Dequeue().([]byte)
				if assertResult == false {
					logger(name, LogWarning, "serial2udpQueue dequeued a non-[]byte node")
				} else {
					//fmt.Printf("UDP out %q\n", toSend)
					//_, err = udpOutputConnection.Write(toSend)
					_, err = udpOutputConnection.Write(toSend)
					if err != nil {
						// TODO do not repeat error for every packet
						// logger(name, LogWarning, err)
						//fmt.Printf("UDP refused for packet %q\n", toSend)
						stats.LostPackets++
					} else {
						//fmt.Printf("UDP sent for packet %q\n", toSend)
					}
				}
			}
			if running == false {
				logger(name, LogInfo, "udpqueue2udp subthread stopped")
				break
			}
			time.Sleep(time.Millisecond * 1)
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
