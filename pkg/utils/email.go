package utils

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"time"
)

type EmailSenderInterface interface {
	SendNotification(to string, content EmailContent) error
}

type EmailSender struct {
	from     string
	password string
	smtpHost string
	smtpPort string
	timeout  time.Duration
}

type EmailContent struct {
	Subject string
	Body    string
}

func NewEmailSender(from, password, smtpHost, smtpPort string) *EmailSender {
	return &EmailSender{
		from:     from,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		timeout:  15 * time.Second,
	}
}

func NewDefaultEmailSender() *EmailSender {
	return NewEmailSender(
		"kudryavtsevasv@mer.ci.nsu.ru",
		"ydzi zydz qysu hzzy",
		"smtp.gmail.com",
		"587",
	)
}

func (e *EmailSender) SendNotification(to string, content EmailContent) error {
	if e.smtpHost == "" || e.smtpPort == "" || e.from == "" {
		return fmt.Errorf("SMTP параметры не настроены")
	}

	msg := []byte(
		"From: " + e.from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + content.Subject + "\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" +
			content.Body + "\r\n",
	)

	conn, err := net.DialTimeout("tcp", e.smtpHost+":"+e.smtpPort, e.timeout)
	if err != nil {
		return fmt.Errorf("SMTP connection failed: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, e.smtpHost)
	if err != nil {
		return fmt.Errorf("SMTP client creation failed: %v", err)
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: e.smtpHost}
		if err = client.StartTLS(config); err != nil {
			return fmt.Errorf("STARTTLS failed: %v", err)
		}
	}

	auth := smtp.PlainAuth("", e.from, e.password, e.smtpHost)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %v", err)
	}

	if err = client.Mail(e.from); err != nil {
		return fmt.Errorf("SMTP MAIL FROM failed: %v", err)
	}

	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("SMTP RCPT TO failed: %v", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA failed: %v", err)
	}

	if _, err = w.Write(msg); err != nil {
		return fmt.Errorf("SMTP message write failed: %v", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("SMTP message close failed: %v", err)
	}

	return client.Quit()
}
