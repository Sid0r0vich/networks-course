package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Ошибка при получении сетевых интерфейсов:", err)
		os.Exit(1)
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println("Ошибка при получении адресов интерфейса:", err)
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			var mask net.IPMask

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				mask = v.Mask
			case *net.IPAddr:
				ip = v.IP
				mask = net.CIDRMask(32, 32)
			default:
				continue
			}

			if ip.IsLoopback() {
				continue
			}

			fmt.Printf("IP-адрес: %s\n", ip.String())
			fmt.Printf("Маска сети: %s\n", mask.String())
		}
	}
}
