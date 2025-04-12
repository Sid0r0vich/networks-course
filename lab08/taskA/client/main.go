package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

const (
	serverAddress   = "localhost:12345"
	packetSize      = 8
	timeoutDuration = 2 * time.Second
	lossProbability = 0.3
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", serverAddress)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	file, err := os.Open("test_file.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, packetSize)
	seqNum := byte(0)

	for {
		n, err := file.Read(buffer[1:])
		if n == 0 {
			break
		}
		if err != nil {
			fmt.Printf("Error reading file:", err)
		}

		buffer[0] = seqNum

		for {
			if rand.Float32() > float32(lossProbability) {
				_, err := conn.Write(buffer[:n+1])
				if err != nil {
					fmt.Println("Error sending packet:", err)
				} else {
					fmt.Printf("Sent packet %d\n", seqNum)
				}
			} else {
				fmt.Printf("Packet %d lost\n", seqNum)
			}

			conn.SetReadDeadline(time.Now().Add(timeoutDuration))

			var ack [1]byte
			n, _, err := conn.ReadFrom(ack[:])
			if n > 0 && ack[0] == seqNum {
				fmt.Printf("Received ACK for packet %d\n", ack[0])
				break
			}

			if err != nil {
				fmt.Println("Timeout, resending packet")
			}
		}

		seqNum ^= 1
	}

	fmt.Println("File transfer complete.")
}
