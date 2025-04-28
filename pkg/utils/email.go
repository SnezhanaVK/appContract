// utils/email.go
package utils

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

// Новый метод для создания отправителя с настройками по умолчанию
func NewDefaultEmailSender() *EmailSender {
    return NewEmailSender(
        "kudryavtsevasv@mer.ci.nsu.ru",
        "ydzi zydz qysu hzzy",
        "smtp.gmail.com",
        "587",
    )
}

type EmailContent struct {
    Subject string
    Body    string
}

func (e *EmailSender) SendNotification(to string, content EmailContent) error {
    if e.smtpHost == "" || e.smtpPort == "" || e.from == "" {
        return fmt.Errorf("SMTP параметры не настроены")
    }
    auth := smtp.PlainAuth("", e.from, e.password, e.smtpHost)

    msg := []byte(
        "To: " + to + "\r\n" +
            "Subject: " + content.Subject + "\r\n" +
            "Content-Type: text/html; charset=UTF-8\r\n" +
            "\r\n" +
            content.Body + "\r\n",
    )

    err := smtp.SendMail(e.smtpHost+":"+e.smtpPort, auth, e.from, []string{to}, msg)
    if err != nil {
        log.Printf("Ошибка отправки письма: %v", err)
        return err
    }
    return nil
}
