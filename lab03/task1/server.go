package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
		return
	}

	requestLine := string(buffer)
	parts := strings.Split(requestLine, "\n")
	if len(parts) < 1 {
		return
	}
	fileName := strings.Split(parts[0], " ")[1][1:]

	content, err := ioutil.ReadFile("../files/" + fileName)
	var response string

	if err != nil {
		response = "HTTP/1.1 404 Not Found\r\n" +
			"Content-Length: 0\r\n\r\n"
	} else {
		response = "HTTP/1.1 200 OK\r\n" +
			fmt.Sprintf("Content-Length: %d\r\n", len(content)) +
			"Content-Type: text/plain\r\n\r\n" +
			string(content)
	}

	conn.Write([]byte(response))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Print("Команда должна содержать порт")
		return
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Print("Неверный порт")
		return
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		os.Exit(1)
	}

	defer listener.Close()
	fmt.Printf("Сервер запущен на порту %d\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка при подключении:", err)
			continue
		}

		handleConnection(conn)
	}
}
