package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Ошибка подключения к серверу:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Введите команду для выполнения на сервере:")
	reader := bufio.NewReader(os.Stdin)
	command, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения команды:", err)
		return
	}

	_, err = conn.Write([]byte(command + "\n"))
	if err != nil {
		fmt.Println("Ошибка отправки команды:", err)
		return
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
		return
	}

	fmt.Println("Результат выполнения команды:")
	fmt.Println(string(buffer[:n]))
}
