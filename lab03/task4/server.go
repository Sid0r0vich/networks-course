package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
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
	if len(os.Args) != 3 {
		fmt.Print("Команда должна содержать порт и максимальное количество одновременных соединений")
		return
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Print("Неверный порт")
		return
	}

	concurrencyLevel, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Print("Неверное максимальное количество одновременных соединений")
		return
	}
	if concurrencyLevel < 0 {
		fmt.Print("Максимальное количество одновременных соединений не может быть меньше нуля")
		return
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		os.Exit(1)
	}

	defer listener.Close()
	fmt.Printf("Сервер запущен на порту %d\n", port)

	var wg sync.WaitGroup
	var tokens = make(chan struct{}, concurrencyLevel)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка при подключении:", err)
			continue
		}

		wg.Add(1)
		tokens <- struct{}{}
		go func(c net.Conn) {
			defer wg.Done()
			handleConnection(conn)
			<-tokens
		}(conn)
	}

	wg.Wait()
}
