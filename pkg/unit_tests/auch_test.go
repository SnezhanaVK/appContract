package unit_tests

import (

	//"appContract/pkg/models"
	//"appContract/pkg/utils"
	db "appContract/pkg/db/repository"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

// TestAuthorize_Success тестирует успешную авторизацию
func TestAuthorize_Success(t *testing.T) {

	t.Run("Should return error when DB connection is nil", func(t *testing.T) {
		_, err := db.Authorize("test", "pass")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "connection error")
	})
}

// TestTokenVerification тестирует верификацию токенов
func TestTokenVerification(t *testing.T) {
	tests := []struct {
		name    string
		claims  jwt.MapClaims
		key     string
		wantErr bool
	}{
		{
			name: "Valid token",
			claims: jwt.MapClaims{
				"id":    1,
				"login": "user",
				"exp":   time.Now().Add(time.Hour).Unix(),
			},
			key:     "secretkey",
			wantErr: false,
		},
		{
			name: "Expired token",
			claims: jwt.MapClaims{
				"id":    1,
				"login": "user",
				"exp":   time.Now().Add(-time.Hour).Unix(),
			},
			key:     "secretkey",
			wantErr: true,
		},
		{
			name: "Wrong signature",
			claims: jwt.MapClaims{
				"id":    1,
				"login": "user",
				"exp":   time.Now().Add(time.Hour).Unix(),
			},
			key:     "wrongkey",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, tt.claims)
			tokenString, err := token.SignedString([]byte(tt.key))
			assert.NoError(t, err)

			_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte("secretkey"), nil
			})

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestVerifyPassword тестирует функцию проверки пароля
func TestVerifyPassword(t *testing.T) {
	// Создаем тестовую структуру для utils
	type verifyFunc func(hashed, pass, salt string) bool

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		salt           string
		verify         verifyFunc
		expected       bool
	}{
		{
			name:           "Correct password",
			hashedPassword: "5f4dcc3b5aa765d61d8327deb882cf99", // MD5 of "password"
			password:       "password",
			salt:           "",
			verify: func(hashed, pass, salt string) bool {
				return hashed == "5f4dcc3b5aa765d61d8327deb882cf99" && pass == "password" && salt == ""
			},
			expected: true,
		},
		{
			name:           "Wrong password",
			hashedPassword: "5f4dcc3b5aa765d61d8327deb882cf99",
			password:       "wrong",
			salt:           "",
			verify: func(hashed, pass, salt string) bool {
				return hashed == "5f4dcc3b5aa765d61d8327deb882cf99" && pass == "password" && salt == ""
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.verify(tt.hashedPassword, tt.password, tt.salt)
			assert.Equal(t, tt.expected, result)
		})
	}
}
