package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "255.255.255.255:9999")
	if err != nil {
		fmt.Println("Ошибка разрешения адреса:", err)
		return
	}

	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		fmt.Println("Ошибка создания сокета:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Сервер запущен. Широковещательная рассылка времени...")

	for {
		currentTime := time.Now().Format("2006-01-02 15:04:05")

		_, err := conn.WriteToUDP([]byte(currentTime), addr)
		if err != nil {
			fmt.Println("Ошибка отправки данных:", err)
		}
		time.Sleep(1 * time.Second)
	}
}
