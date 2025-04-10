package utils

import (
	"net/smtp"
)

// отправляет сообщение на SMTP-сервер
func SendEmail(to string, subject string, body string) error {

	auth := smtp.PlainAuth("", "your_email@gmail.com", "your_password", "smtp.gmail.com")

	
	msg := []byte("To: " + to + "\r\nSubject: " + subject + "\r\n\r\n" + body)

	
	err := smtp.SendMail("smtp.gmail.com:587", auth, "your_email@gmail.com", []string{to}, msg)
	if err != nil {
		return err
	}

	return nil
}

