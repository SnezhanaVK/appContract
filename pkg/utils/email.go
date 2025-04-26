package utils

// email.go в пакете utils
import (
	"fmt"
	"log"
	"net/smtp"
)

type EmailSender struct {
	from     string
	password string
	smtpHost string
	smtpPort string
}

func NewEmailSender(from, password, smtpHost, smtpPort string) *EmailSender {
	return &EmailSender{
		from:     from,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

type EmailContent struct {
	Subject string
	Body    string
}

func (e *EmailSender) SendNotification(to string, content EmailContent) error {
    auth := smtp.PlainAuth("", e.from, e.password, e.smtpHost)

    msg := []byte(
        "To: " + to + "\r\n" +
            "Subject: " + content.Subject + "\r\n" +
            "Content-Type: text/html; charset=UTF-8\r\n" +
            "\r\n" +
            content.Body + "\r\n",
    )

    // Таймаут 10 секунд для SMTP соединения
    err := smtp.SendMail(e.smtpHost+":"+e.smtpPort, auth, e.from, []string{to}, msg)
    if err != nil {
        log.Printf("SMTP ошибка: %v", err)
        return fmt.Errorf("ошибка отправки email: %v", err)
    }
    
    return nil
}

