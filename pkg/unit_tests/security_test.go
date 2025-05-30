package unit_tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"appContract/pkg/db"
	// repository "appContract/pkg/db/repository"
	"appContract/pkg/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Глобальная переменная для хранения экземпляра пула базы данных
var dbPool *pgxpool.Pool

// Функция для установки связи с базой данных перед каждым набором тестов
func TestMain(m *testing.M) {
    // Сначала устанавливаем соединение с базой данных
    db.ConnectDB()
    defer db.CloseDB()

    // Выполняем сами тесты
    code := m.Run()

    // Завершаем выполнение программы с результатом выполнения тестов
    if code > 0 {
        panic("Tests failed.")
    }
}



func TestSQLInjectionInLogin(t *testing.T) {
    // Тестирование входа с SQL-инъекцией
    payload := `{
        "login": "admin'--",
        "password": "anything"
    }`
    req, err := http.NewRequest("POST", "/api/authorizations", strings.NewReader(payload))
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(handlers.Login)

    handler.ServeHTTP(rr, req)

    // Проверка статуса HTTP-запроса
    if status := rr.Code; status == http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusUnauthorized)
    }

    // Проверка наличия технической информации об ошибке
    if strings.Contains(rr.Body.String(), "SQL syntax") {
        t.Error("handler returned SQL error details, potential security risk")
    }
}


func TestSensitiveDataExposure(t *testing.T) {
    // Проверяем отсутствие утечки чувствительных данных в API
    req, err := http.NewRequest("GET", "/api/users/1", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(handlers.GetUserID)

    handler.ServeHTTP(rr, req)

    // Проверяем, что ответ не содержит конфиденциальных полей
    responseBody := rr.Body.String()
    sensitiveFields := []string{
        "password_hash",
        "salt",
        "password_algorithm",
    }

    for _, field := range sensitiveFields {
        if strings.Contains(responseBody, field) {
            t.Errorf("sensitive field '%s' exposed in API response", field)
        }
    }

    // Проверяем отсутствие пароля в хэшированной форме
    if strings.Contains(responseBody, "$2a$") { // Префикс Bcrypt-хэшей
        t.Error("password hash exposed in API response")
    }
}