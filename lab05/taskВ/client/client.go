package main

import (
	"fmt"
	"net"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":9999")
	if err != nil {
		fmt.Println("Ошибка разрешения адреса:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Ошибка создания сокета:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Клиент запущен. Ожидание сообщений...")

	buffer := make([]byte, 1024)

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Ошибка чтения данных:", err)
			continue
		}
		fmt.Println("Полученное время:", string(buffer[:n]))
	}
}
