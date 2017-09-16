package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/jacobsa/go-serial/serial"
	"github.com/oleiade/lane"
)

// UDPSerialThread : start a loop for a specified port
func UDPSerialThread(name string, portConfig PortConfig, stopChannel chan string, killChannel chan bool) {
	defer func() { stopChannel <- portConfig.Name }()

	logger(name, LogInfo, "Starting thread")

	var running = true

	// Get serial port TTY
	ttyName, err := getPortTTY(definitions, portConfig.Name)
	if err != nil {
		logger(name, LogError, err)
	}

	// Serial port configuration
	serialPortOptions := serial.OpenOptions{
		PortName:        ttyName,
		BaudRate:        portConfig.BaudRate,
		DataBits:        portConfig.DataBits,
		StopBits:        portConfig.StopBits,
		MinimumReadSize: 1,
	}

	// UDP Input Address
	udpInputAddress, err := net.ResolveUDPAddr("udp", portConfig.UDPInputIP+":"+strconv.Itoa(portConfig.UDPInputPort))
	if err != nil {
		logger(name, LogError, err)
		return
	}

	// UDP Output Address
	udpOutputAddress, err := net.ResolveUDPAddr("udp", portConfig.UDPOutputIP+":"+strconv.Itoa(portConfig.UDPOutputPort))
	if err != nil {
		logger(name, LogError, err)
		return
	}

	// Open serial port
	serialPort, err := serial.Open(serialPortOptions)
	if err != nil {
		logger(name, LogError, err)
		return
	}
	logger(name, LogInfo, "Opened "+serialPortOptions.PortName)
	defer serialPort.Close()

	// Open UDP input connection
	udpInputConnection, err := net.ListenUDP("udp", udpInputAddress)
	if err != nil {
		logger(name, LogError, err)
		return
	}
	logger(name, LogInfo, "Listening on "+udpInputAddress.String())
	defer udpInputConnection.Close()

	// Open UDP output connection
	udpOutputConnection, err := net.DialUDP("udp", nil, udpOutputAddress)
	if err != nil {
		logger(name, LogError, err)
		return
	}
	logger(name, LogInfo, "Sending to "+udpOutputAddress.String())
	defer udpOutputConnection.Close()

	var udpBuffer = make([]byte, 1024)
	var serialBuffer = make([]byte, 1024)

	var serial2udpQueue = lane.NewQueue()

	var internalWaitGroup sync.WaitGroup

	internalWaitGroup.Add(1)
	go func() {
		defer internalWaitGroup.Done()
		var readLength int
		var readAddr *net.UDPAddr
		for {
			readLength, readAddr, err = udpInputConnection.ReadFromUDP(udpBuffer)
			if err != nil {
				logger(name, LogWarning, err)
			} else {
				fmt.Println("UDP: ", string(udpBuffer[0:readLength]), " from ", readAddr)
				go func() {
					_, err = serialPort.Write(udpBuffer[:readLength])
					if err != nil {
						// TODO do not repeat error for every packet
						logger(name, LogWarning, err)
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
		for {
			readLength, err = serialPort.Read(serialBuffer)
			if err == io.EOF {
				logger(name, LogWarning, err)
			} else {
				fmt.Printf("Serial: %q\n", serialBuffer[:readLength])
				serial2udpQueue.Enqueue(serialBuffer[:readLength])
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
				fmt.Printf("queued=%d\n", serial2udpQueue.Size())
				var toSend, assertResult = serial2udpQueue.Dequeue().([]byte)
				if assertResult == false {
					logger(name, LogWarning, "serial2udpQueue dequeued a non-[]byte node")
				} else {
					//_, err = udpOutputConnection.Write(toSend)
					_, _, err = udpOutputConnection.WriteMsgUDP(toSend, nil, nil)
					if err != nil {
						// TODO do not repeat error for every packet
						logger(name, LogWarning, err)
						fmt.Printf("Refused for packet %q\n", toSend)
					} else {
						fmt.Printf("UDP sent for packet %q\n", toSend)
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
