package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
)

func checkPort(ip string, port int, wg *sync.WaitGroup, availablePorts chan<- int) {
	defer wg.Done()
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		availablePorts <- port
		return
	}
	defer conn.Close()
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Использование: go run main.go <IP-адрес> <начальный порт> <конечный порт>")
		return
	}

	ip := os.Args[1]
	startPort, err := strconv.Atoi(os.Args[2])
	if err != nil || startPort < 0 || startPort > 65535 {
		fmt.Println("Некорректный начальный порт. Убедитесь, что это число от 0 до 65535.")
		return
	}

	endPort, err := strconv.Atoi(os.Args[3])
	if err != nil || endPort < startPort || endPort > 65535 {
		fmt.Println("Некорректный конечный порт. Убедитесь, что это число от начального порта до 65535.")
		return
	}

	var wg sync.WaitGroup
	availablePorts := make(chan int)

	go func() {
		for port := range availablePorts {
			fmt.Printf("Порт %d доступен\n", port)
		}
	}()

	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go checkPort(ip, port, &wg, availablePorts)
	}

	wg.Wait()
	close(availablePorts)
}
