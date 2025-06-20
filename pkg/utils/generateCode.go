package utils

//
import (
	"fmt"
	"math/rand"
	"time"
)


func GenerateVerificationCode() string {
	
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	code := random.Intn(100000)
	return fmt.Sprintf("%05d", code)
}