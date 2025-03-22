package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

const (
	smtpServer = "smtp.gmail.com:465"
	username   = "dmitriysidorcool@gmail.com"
	password   = "tall fenz eowh vkdb"
)

func sendEmail(to string, subject string, body string, imagePath string) error {
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

	imageData, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return fmt.Errorf("ошибка чтения изображения: %v", err)
	}

	encodedImage := base64.StdEncoding.EncodeToString(imageData)

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=\"boundary\"\r\n\r\n", username, to, subject))
	builder.WriteString("--boundary\r\n")
	builder.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n")
	builder.WriteString(body + "\r\n\r\n")
	builder.WriteString("--boundary\r\n")
	builder.WriteString("Content-Type: image/jpeg; name=\"" + imagePath + "\"\r\n")
	builder.WriteString("Content-Transfer-Encoding: base64\r\n")
	builder.WriteString("Content-Disposition: attachment; filename=\"" + imagePath + "\"\r\n\r\n")
	builder.WriteString(encodedImage + "\r\n\r\n")
	builder.WriteString("--boundary--\r\n.")

	if err := sendCommand(conn, builder.String()); err != nil {
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
	if len(os.Args) != 3 {
		fmt.Errorf("команда должна содержать адрес получателя и путь до изображения")
		return
	}
	recipient := os.Args[1]
	imagePath := os.Args[2]
	subject := "Тестовое сообщение"
	body := "Это тестовое сообщение"

	if err := sendEmail(recipient, subject, body, imagePath); err != nil {
		fmt.Println("Ошибка при отправке сообщения:", err)
	} else {
		fmt.Println("Сообщение отправлено успешно!")
	}
}
