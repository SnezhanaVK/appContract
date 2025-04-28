package service

// //autch_server.go в папке servise
// import (
// 	"appContract/pkg/models"
// 	"appContract/pkg/utils"
// 	"fmt"
// 	"log"
// 	"sync"
// 	"time"
// )
// type codeStorage struct {
// 	mu    sync.Mutex
// 	codes map[string]codeEntry
// }

// type codeEntry struct {
// 	code      string
// 	createdAt time.Time
// }

// var storage = &codeStorage{
// 	codes: make(map[string]codeEntry),
// }
// // Добавляем код в хранилище
// func (s *codeStorage) addCode(email, code string) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	s.codes[email] = codeEntry{
// 		code:      code,
// 		createdAt: time.Now(),
// 	}
// }

// // Проверяем код
// func (s *codeStorage) verifyCode(email, code string) bool {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	log.Printf("Checking code for %s, input code: %s", email, code)
//     log.Printf("Stored codes: %+v", s.codes)

// 	entry, exists := s.codes[email]
// 	if !exists {
// 		log.Println("No code found for email")
// 		return false
// 	}

// 	// Удаляем код после проверки (одноразовый)
// 	delete(s.codes, email)

// 	// Проверяем срок действия (10 минут)
// 	if time.Since(entry.createdAt) > 10*time.Minute {
// 		return false
// 	}

// 	return entry.code == code
// }

// // Запускаем очистку старых кодов
// func init() {
// 	go func() {
// 		for {
// 			time.Sleep(1 * time.Hour)
// 			storage.mu.Lock()
// 			for email, entry := range storage.codes {
// 				if time.Since(entry.createdAt) > 10*time.Minute {
// 					delete(storage.codes, email)
// 				}
// 			}
// 			storage.mu.Unlock()
// 		}
// 	}()
// }
// func SendingCode(user models.Users) (string, error) {
// 	// Генерируем случайный код
// 	verificationCode := utils.GenerateVerificationCode()

// 	// Проверяем наличие email у пользователя
// 	if user.Email == "" {
// 		return "", fmt.Errorf("user email is empty")
// 	}

// 	// Формируем содержимое письма
// 	emailContent := utils.EmailContent{
// 		Subject: "Ваш код подтверждения",
// 		Body:    fmt.Sprintf("Ваш код подтверждения: <strong>%s</strong>", verificationCode),
// 	}

// 	// Отправляем код пользователю
// 	err := emailSender.SendNotification(user.Email, emailContent)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to send verification code: %w", err)
// 	}
//  // Сохраняем код в памяти
//  	storage.addCode(user.Email, verificationCode)
// 	return verificationCode, nil
// }

// func VerifyCode(email, code string) bool {
// 	return storage.verifyCode(email, code)
// }




