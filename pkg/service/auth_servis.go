package service

import (
	"appContract/pkg/models"
	"appContract/pkg/utils"
	"fmt"
	"log"
	"sync"
	"time"
)

type codeStorage struct {
	mu    sync.Mutex
	codes map[string]codeEntry
}

type codeEntry struct {
	code      string
	createdAt time.Time
}

var storage = &codeStorage{
	codes: make(map[string]codeEntry),
}

func (s *codeStorage) addCode(email, code string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.codes[email] = codeEntry{
		code:      code,
		createdAt: time.Now(),
	}
}

func (s *codeStorage) verifyCode(email, code string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Printf("Checking code for %s, input code: %s", email, code)
	log.Printf("Stored codes: %+v", s.codes)

	entry, exists := s.codes[email]
	if !exists {
		log.Println("No code found for email")
		return false
	}

	delete(s.codes, email)

	if time.Since(entry.createdAt) > 10*time.Minute {
		return false
	}

	return entry.code == code
}

func init() {
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			storage.mu.Lock()
			for email, entry := range storage.codes {
				if time.Since(entry.createdAt) > 10*time.Minute {
					delete(storage.codes, email)
				}
			}
			storage.mu.Unlock()
		}
	}()
}
func SendingCode(user models.Users) (string, error) {
	verificationCode := utils.GenerateVerificationCode()

	if user.Email == "" {
		return "", fmt.Errorf("user email is empty")
	}

	emailContent := utils.EmailContent{
    Subject: "Ваш код подтверждения",
    Body: fmt.Sprintf(`
        <div style="
            background-color: rgb(67, 73, 72);
            color: #ffffff;
            padding: 25px;
            font-family: 'Segoe UI', Arial, sans-serif;
            border-radius: 8px;
            max-width: 600px;
            margin: 0 auto;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        ">
            <h2 style="color: #ffffff; font-size: 18px; margin-top: 0;">Код подтверждения</h2>
            <p style="font-size: 15px; line-height: 1.5;">Уважаемый пользователь,</p>
            <p style="font-size: 15px; line-height: 1.5;">
                Ваш код подтверждения: <strong style="
                    color: #ffffff;
                    font-size: 20px;
                    letter-spacing: 2px;
                    background: rgba(255,255,255,0.1);
                    padding: 8px 12px;
                    border-radius: 4px;
                    display: inline-block;
                ">%s</strong>
            </p>
            <p style="font-size: 14px; line-height: 1.5; color: #cccccc;">
                Никому не сообщайте этот код.
            </p>
            <div style="margin-top: 20px; padding-top: 10px; border-top: 1px solid rgba(255, 255, 255, 0.1);">
                <p style="font-size: 13px; color: #cccccc;">Это автоматическое уведомление. Пожалуйста, не отвечайте на него.</p>
            </div>
        </div>
    `, verificationCode),
}

	err := emailSender.SendNotification(user.Email, emailContent)
	if err != nil {
		return "", fmt.Errorf("failed to send verification code: %w", err)
	}
	storage.addCode(user.Email, verificationCode)
	return verificationCode, nil
}

func VerifyCode(email, code string) bool {
	return storage.verifyCode(email, code)
}
