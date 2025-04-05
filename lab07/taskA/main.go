package main

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

func main() {
	address := net.UDPAddr{
		Port: 12345,
		IP:   net.ParseIP("127.0.0.1"),
	}

	conn, err := net.ListenUDP("udp", &address)
	if err != nil {
		fmt.Println("Ошибка при создании сокета:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Сервер запущен на %s:%d. Ожидание пакетов...\n", address.IP, address.Port)

	rand.Seed(time.Now().UnixNano())

	for {
		buffer := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Ошибка при получении данных:", err)
			continue
		}

		message := string(buffer[:n])
		fmt.Printf("Получен пакет от %s: %s\n", addr.String(), message)

		if rand.Float32() < 0.2 {
			fmt.Println("Пакет потерян.")
			continue
		}

		response := strings.ToUpper(message)

		if _, err := conn.WriteToUDP([]byte(response), addr); err != nil {
			fmt.Println("Ошибка при отправке ответа:", err)
			continue
		}

		fmt.Printf("Отправлен ответ: %s на %s\n", response, addr.String())
	}
}
