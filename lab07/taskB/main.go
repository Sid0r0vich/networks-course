package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	local, err := net.ResolveUDPAddr("udp", "localhost:12344")
	if err != nil {
		panic(err)
	}
	remote, err := net.ResolveUDPAddr("udp", "localhost:12345")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", local, remote)
	if err != nil {
		fmt.Println("Ошибка при подключении к серверу:", err)
		return
	}
	defer conn.Close()

	for i := 1; i <= 10; i++ {
		startTime := time.Now()
		message := fmt.Sprintf("Ping %d %v", i, startTime.Format("15:04:05"))

		if _, err := conn.Write([]byte(message)); err != nil {
			fmt.Println("Ошибка при отправке сообщения:", err)
			continue
		}

		conn.SetReadDeadline(time.Now().Add(1 * time.Second))

		buffer := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buffer)

		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("Request timed out")
			} else {
				fmt.Println("Ошибка при получении ответа:", err)
			}
			continue
		}

		rtt := time.Since(startTime).Seconds()
		responseMessage := string(buffer[:n])

		fmt.Printf("message: %s\nRTT: %.3f\n", responseMessage, rtt)

		time.Sleep(500 * time.Millisecond)
	}
}
