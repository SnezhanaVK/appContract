package utils

import (
	"net/smtp"
)

// SendEmail отправляет сообщение на SMTP-сервер
func SendEmail(to string, subject string, body string) error {
	// Создаем новый экземпляр структуры smtp.Auth
	auth := smtp.PlainAuth("", "your_email@gmail.com", "your_password", "smtp.gmail.com")

	// Создаем новый экземпляр структуры []byte
	msg := []byte("To: " + to + "\r\nSubject: " + subject + "\r\n\r\n" + body)

	// Отправляем сообщение на SMTP-сервер
	err := smtp.SendMail("smtp.gmail.com:587", auth, "your_email@gmail.com", []string{to}, msg)
	if err != nil {
		return err
	}

	return nil
}

