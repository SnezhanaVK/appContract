// utils/password.go
package utils

import (
	"crypto/rand"
	"math/big"
)

// GenerateStrongPassword создает надежный постоянный пароль
func GenerateStrongPassword() (string, error) {
	const (
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits    = "0123456789"
		special   = "!@#$%^&*()_+-="
		allChars  = lowercase + uppercase + digits + special
	)

	// Минимальные требования: 3 lowercase, 2 uppercase, 2 digits, 1 special
	password := make([]byte, 16) // Длина 16 символов

	// Заполняем минимальные требования
	password[0] = lowercase[randInt(len(lowercase))]
	password[1] = lowercase[randInt(len(lowercase))]
	password[2] = lowercase[randInt(len(lowercase))]
	password[3] = uppercase[randInt(len(uppercase))]
	password[4] = uppercase[randInt(len(uppercase))]
	password[5] = digits[randInt(len(digits))]
	password[6] = digits[randInt(len(digits))]
	password[7] = special[randInt(len(special))]

	// Заполняем оставшиеся позиции
	for i := 8; i < len(password); i++ {
		password[i] = allChars[randInt(len(allChars))]
	}

	// Перемешиваем пароль
	shufflePassword(password)

	return string(password), nil
}

func randInt(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}

func shufflePassword(p []byte) {
	for i := len(p) - 1; i > 0; i-- {
		j := randInt(i + 1)
		p[i], p[j] = p[j], p[i]
	}
}
