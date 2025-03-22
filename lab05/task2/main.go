package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	smtpServer = "smtp.gmail.com:465"
	username   = "dmitriysidorcool@gmail.com"
	password   = "tall fenz eowh vkdb"
)

func sendEmail(to string, subject string, body string) error {
	// Устанавливаем TLS-соединение
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "smtp.gmail.com",
	}

	conn, err := tls.Dial("tcp", smtpServer, tlsConfig)
	if err != nil {
		return fmt.Errorf("ошибка подключения: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	if _, err := reader.ReadString('\n'); err != nil {
		return fmt.Errorf("ошибка чтения приветственного сообщения: %v", err)
	}

	if err := sendCommand(conn, "EHLO localhost"); err != nil {
		return err
	}

	if err := sendCommand(conn, "AUTH LOGIN"); err != nil {
		return err
	}

	if err := sendCommand(conn, base64.StdEncoding.EncodeToString([]byte(username))); err != nil {
		return err
	}
	if err := sendCommand(conn, base64.StdEncoding.EncodeToString([]byte(password))); err != nil {
		return err
	}

	if err := sendCommand(conn, "MAIL FROM:<"+username+">"); err != nil {
		return err
	}

	if err := sendCommand(conn, "RCPT TO:<"+to+">"); err != nil {
		return err
	}

	if err := sendCommand(conn, "DATA"); err != nil {
		return err
	}

	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n.", username, to, subject, body)
	if err := sendCommand(conn, message); err != nil {
		return err
	}

	if err := sendCommand(conn, "QUIT"); err != nil {
		return err
	}

	return nil
}

func sendCommand(conn net.Conn, command string) error {
	_, err := fmt.Fprintf(conn, command+"\r\n")
	if err != nil {
		return fmt.Errorf("ошибка отправки команды: %v", err)
	}

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	if !strings.HasPrefix(response, "250") && !strings.HasPrefix(response, "354") && !strings.HasPrefix(response, "334") && !strings.HasPrefix(response, "235") && !strings.HasPrefix(response, "221") {
		return fmt.Errorf("неожиданный ответ от сервера: %s", response)
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Errorf("команда должна содержать адрес получателя")
		return
	}
	recipient := os.Args[1]
	subject := "Тестовое сообщение"
	body := "Это тестовое сообщение"

	if err := sendEmail(recipient, subject, body); err != nil {
		fmt.Println("Ошибка при отправке сообщения:", err)
	} else {
		fmt.Println("Сообщение отправлено успешно!")
	}
}
