package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// GenerateSalt создает криптографически безопасную случайную соль заданной длины
func GenerateSalt(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("salt length must be positive")
	}

	// Минимальная рекомендуемая длина соли - 16 байт
	if length < 16 {
		length = 16
	}

	// Создаем буфер для случайных байтов
	saltBytes := make([]byte, length)

	// Заполняем буфер криптографически безопасными случайными байтами
	_, err := rand.Read(saltBytes)
	if err != nil {
		return "", err
	}

	// Кодируем в base64 для удобного хранения в БД
	salt := base64.URLEncoding.EncodeToString(saltBytes)

	// Обрезаем до нужной длины (base64 увеличивает длину)
	if len(salt) > length {
		salt = salt[:length]
	}

	return salt, nil
}

// HashPassword создает bcrypt хеш из пароля и соли
func HashPassword(password, salt string) (string, error) {
	// Комбинируем пароль и соль перед хешированием
	combined := password + salt

	// Генерируем хеш с стандартной стоимостью (можно увеличить для большей безопасности)
	hashed, err := bcrypt.GenerateFromPassword([]byte(combined), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

// VerifyPassword проверяет соответствие пароля хешу
func VerifyPassword(hashedPassword, password, salt string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password+salt),
	)
	return err == nil
}
