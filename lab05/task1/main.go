package main

import (
	"fmt"
	"net/smtp"
	"strings"
)

const (
	smtpServer = "smtp.gmail.com"
	smtpPort   = "587"
	username   = "dmitriysidorcool@gmail.com"
	password   = "tall fenz eowh vkdb"
	sender     = "dmitriysidorcool@gmail.com"
)

func sendEmailText(to string, subject string, body string) error {
	headers := make(map[string]string)
	headers["Subject"] = subject

	var msg strings.Builder
	for key, value := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	auth := smtp.PlainAuth("", username, password, smtpServer)
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, sender, []string{to}, []byte(msg.String()))
	return err
}

func sendEmailHTML(to string, subject string, html string) error {
	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	msg := "Subject: " + subject + "\n" + headers + "\n\n" + html

	auth := smtp.PlainAuth("", username, password, smtpServer)
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, sender, []string{to}, []byte(msg))
	return err
}

func main() {
	recipient := "dmitriysidorclash@gmail.com"

	txtSubject := "Тестовое текстовое сообщение"
	txtBody := "Это текстовое сообщение."
	if err := sendEmailText(recipient, txtSubject, txtBody); err != nil {
		fmt.Println("Ошибка при отправке текстового сообщения:", err)
	} else {
		fmt.Println("Текстовое сообщение отправлено успешно!")
	}

	htmlSubject := "Тестовое HTML сообщение"
	htmlBody := "<h1>Это HTML сообщение</h1><p>Содержимое HTML сообщения.</p>"
	if err := sendEmailHTML(recipient, htmlSubject, htmlBody); err != nil {
		fmt.Println("Ошибка при отправке HTML сообщения:", err)
	} else {
		fmt.Println("HTML сообщение отправлено успешно!")
	}
}
