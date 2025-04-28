package utils

//
import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateVerificationCode создает 6-значный цифровой код для подтверждения почты
// Возвращает строку с ведущими нулями (пример: "043789")
func GenerateVerificationCode() string {
	// Инициализация генератора с уникальным сидом
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	
	// Генерация числа в диапазоне 000000-999999
	code := random.Intn(100000)
	
	// Форматирование с ведущими нулями
	return fmt.Sprintf("%05d", code)
}