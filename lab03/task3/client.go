package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Использование: <хост> <порт> <имя файла>")
		return
	}

	serverHost := os.Args[1]
	serverPort, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Неверный порт")
		return
	}
	fileName := os.Args[3]

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverHost, serverPort))
	if err != nil {
		fmt.Printf("Ошибка создания соединения с сервером:", err)
		return
	}
	defer conn.Close()

	httpRequest := fmt.Sprintf("GET /%s HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n", fileName, serverHost)

	_, err = conn.Write([]byte(httpRequest))
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return
	}

	response, err := ioutil.ReadAll(conn)
	if err != nil {
		fmt.Println("Ошибка чтения запроса:", err)
		return
	}

	fmt.Println(string(response))
}
