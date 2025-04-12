package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
)

const (
	port            = ":12345"
	lossProbability = 0.3
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	file, err := os.Create("received_file.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 1025)
	expectedSeqNum := byte(0)

	for {
		n, clientAddr, err := conn.ReadFrom(buffer)
		if err != nil || n == 0 {
			fmt.Println("Error reading from connection:", err)
			continue
		}

		seqNum := buffer[0]
		data := buffer[1:n]

		if seqNum == expectedSeqNum {
			// fmt.Printf("DATA: %s\n", string(buffer[1:n]))
			fmt.Printf("Received packet %d\n", seqNum)
			_, err = file.Write(data)
			if err != nil {
				fmt.Println("Error writing to file:", err)
			}

			expectedSeqNum ^= 1
		}

		if rand.Float32() > float32(lossProbability) {
			ackPacket := []byte{seqNum}
			_, err = conn.WriteTo(ackPacket, clientAddr)
			if err != nil {
				fmt.Println("Error writing to connection:", err)
			} else {
				fmt.Printf("Sent ACK for packet %d\n", seqNum)
			}
		} else {
			fmt.Printf("ACK for packet %d lost\n", seqNum)
		}
	}
}
