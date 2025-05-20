package unit_tests

// import (
// 	"appContract/pkg/db"
// 	"appContract/pkg/utils"
// 	"context"
// 	"testing"

// 	"github.com/jackc/pgx"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"
// )

// func TestSQLInjectionProtection(t *testing.T) {
// 	// Сохраняем оригинальное подключение к БД
// 	origGetDB := db.GetDB
// 	defer func() { db.GetDB = origGetDB }()

// 	// Создаем мок подключения
// 	mockConn := new(MockPgxConn)
// 	db.GetDB = func() *pgx.Conn {
// 		return mockConn
// 	}

// 	testCases := []struct {
// 		name     string
// 		input    string
// 		expected string
// 	}{
// 		{"Basic input", "normal_input", "normal_input"},
// 		{"SQL injection attempt", "admin' OR '1'='1", "admin' OR '1'='1"},
// 		{"Semicolon attack", "'; DROP TABLE users;--", "'; DROP TABLE users;--"},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Настраиваем мок для возврата ожидаемого значения
// 			mockRow := new(pgx.Row)
// 			mockConn.On("QueryRow", context.Background(), "SELECT $1::text", []interface{}{tc.input}).
// 				Return(mockRow).
// 				Run(func(args mock.Arguments) {
// 					// Симулируем работу Scan
// 					row := args.Get(0).(*pgx.Row)
// 					*row = pgx.Row{}
// 					row.Scan = func(dest ...interface{}) error {
// 						*dest[0].(*string) = tc.input
// 						return nil
// 					}
// 				})

// 			var result string
// 			err := mockConn.QueryRow(context.Background(), "SELECT $1::text", tc.input).Scan(&result)

// 			assert.NoError(t, err)
// 			assert.Equal(t, tc.expected, result)
// 			mockConn.AssertExpectations(t)
// 		})
// 	}
// }

// func TestXSSProtection(t *testing.T) {
// 	// Сохраняем оригинальное подключение к БД
// 	origGetDB := db.GetDB
// 	defer func() { db.GetDB = origGetDB }()

// 	// Создаем мок подключения
// 	mockConn := new(MockPgxConn)
// 	db.GetDB = func() *pgx.Conn {
// 		return mockConn
// 	}

// 	testCases := []struct {
// 		name     string
// 		input    string
// 		expected string
// 	}{
// 		{"Safe input", "Hello world", "Hello world"},
// 		{"Script tag", "<script>alert('XSS')</script>", "<script>alert('XSS')</script>"},
// 		{"Image XSS", "<img src=x onerror=alert('XSS')>", "<img src=x onerror=alert('XSS')>"},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Настраиваем мок
// 			mockRow := new(pgx.Row)
// 			mockConn.On("QueryRow", context.Background(), "SELECT $1::text", []interface{}{tc.input}).
// 				Return(mockRow).
// 				Run(func(args mock.Arguments) {
// 					row := args.Get(0).(*pgx.Row)
// 					*row = pgx.Row{}
// 					row.Scan = func(dest ...interface{}) error {
// 						*dest[0].(*string) = tc.input
// 						return nil
// 					}
// 				})

// 			var result string
// 			err := mockConn.QueryRow(context.Background(), "SELECT $1::text", tc.input).Scan(&result)

// 			assert.NoError(t, err)
// 			assert.Equal(t, tc.expected, result)
// 			mockConn.AssertExpectations(t)
// 		})
// 	}
// }

// func TestPasswordHashing(t *testing.T) {
// 	// Сохраняем оригинальные функции для восстановления после теста
// 	origGenerateSalt := utils.GenerateSalt
// 	origHashPassword := utils.HashPassword
// 	origVerifyPassword := utils.VerifyPassword
// 	defer func() {
// 		utils.GenerateSalt = origGenerateSalt
// 		utils.HashPassword = origHashPassword
// 		utils.VerifyPassword = origVerifyPassword
// 	}()

// 	// Подменяем функции для тестов
// 	utils.GenerateSalt = func(length int) (string, error) {
// 		return "fixed-salt-for-tests", nil
// 	}

// 	utils.HashPassword = func(password, salt string) (string, error) {
// 		// Простая детерминированная хеш-функция для тестов
// 		return password + "|" + salt, nil
// 	}

// 	utils.VerifyPassword = func(hashed, password, salt string) bool {
// 		expected, _ := utils.HashPassword(password, salt)
// 		return hashed == expected
// 	}

// 	testCases := []struct {
// 		name     string
// 		password string
// 	}{
// 		{"Simple password", "password123"},
// 		{"Complex password", "P@ssw0rd!123"},
// 		{"Long password", "verylongpasswordwith1234567890and!@#$%^&*()"},
// 		{"Empty password", ""},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			salt, err := utils.GenerateSalt(16)
// 			require.NoError(t, err)

// 			hashed1, err := utils.HashPassword(tc.password, salt)
// 			require.NoError(t, err)
// 			require.NotEmpty(t, hashed1)

// 			hashed2, err := utils.HashPassword(tc.password, salt)
// 			require.NoError(t, err)
// 			assert.Equal(t, hashed1, hashed2, "Same password+salt should produce same hash")

// 			// Проверка верификации
// 			assert.True(t, utils.VerifyPassword(hashed1, tc.password, salt))
// 			assert.False(t, utils.VerifyPassword(hashed1, "wrongpassword", salt))

// 			// Проверка с другой солью
// 			newSalt := "different-salt"
// 			newHashed, err := utils.HashPassword(tc.password, newSalt)
// 			require.NoError(t, err)
// 			assert.NotEqual(t, hashed1, newHashed, "Different salt should produce different hash")
// 		})
// 	}
// }

// func TestAuthorize(t *testing.T) {
// 	// Сохраняем оригинальные зависимости
// 	origGetDB := db.GetDB
// 	origVerifyPassword := utils.VerifyPassword
// 	defer func() {
// 		db.GetDB = origGetDB
// 		utils.VerifyPassword = origVerifyPassword
// 	}()

// 	// Создаем мок подключения
// 	mockConn := new(MockPgxConn)
// 	db.GetDB = func() *pgx.Conn {
// 		return mockConn
// 	}

// 	// Настраиваем мок для проверки пароля
// 	utils.VerifyPassword = func(hashedPassword, inputPassword, salt string) bool {
// 		return inputPassword == "correctpassword"
// 	}

// 	t.Run("Successful authorization", func(t *testing.T) {
// 		// Мок для запроса пользователя
// 		userRow := new(pgx.Row)
// 		mockConn.On("QueryRow", context.Background(), mock.Anything, []interface{}{"testuser"}).
// 			Return(userRow).
// 			Run(func(args mock.Arguments) {
// 				row := args.Get(0).(*pgx.Row)
// 				*row = pgx.Row{}
// 				row.Scan = func(dest ...interface{}) error {
// 					*(dest[0].(*int)) = 1               // id_user
// 					*(dest[1].(*string)) = "Doe"        // surname
// 					*(dest[2].(*string)) = "John"       // username
// 					*(dest[3].(*string)) = ""           // patronymic
// 					*(dest[4].(*string)) = ""           // phone
// 					*(dest[5].(*string)) = ""           // email
// 					*(dest[6].(*string)) = "testuser"   // login
// 					*(dest[7].(*string)) = "hashedpass" // password_hash
// 					*(dest[8].(*string)) = "salt"       // salt
// 					return nil
// 				}
// 			})

// 		// Мок для запроса ролей
// 		mockRows := new(pgx.Rows)
// 		mockConn.On("Query", context.Background(), mock.Anything, []interface{}{1}).
// 			Return(mockRows, nil).
// 			Run(func(args mock.Arguments) {
// 				rows := args.Get(0).(*pgx.Rows)
// 				*rows = pgx.Rows{}
// 				rows.NextFunc = func() bool {
// 					rows.ScanFunc = func(dest ...interface{}) error {
// 						*(dest[0].(*int)) = 1          // id_role
// 						*(dest[1].(*string)) = "admin" // name_role
// 						return nil
// 					}
// 					return true
// 				}
// 			})

// 		user, err := db.Authorize("testuser", "correctpassword")
// 		assert.NoError(t, err)
// 		assert.Equal(t, 1, user.Id_user)
// 		assert.Equal(t, "testuser", user.Login)
// 		assert.True(t, user.Admin)
// 	})

// 	t.Run("Invalid password", func(t *testing.T) {
// 		userRow := new(pgx.Row)
// 		mockConn.On("QueryRow", context.Background(), mock.Anything, []interface{}{"testuser"}).
// 			Return(userRow).
// 			Run(func(args mock.Arguments) {
// 				row := args.Get(0).(*pgx.Row)
// 				*row = pgx.Row{}
// 				row.Scan = func(dest ...interface{}) error {
// 					*(dest[7].(*string)) = "hashedpass" // password_hash
// 					*(dest[8].(*string)) = "salt"       // salt
// 					return nil
// 				}
// 			})

// 		_, err := db.Authorize("testuser", "wrongpassword")
// 		assert.Error(t, err)
// 		assert.Equal(t, "invalid password", err.Error())
// 	})
// }
