package unit_tests

import (
	"appContract/pkg/middleware"
	"appContract/pkg/routers"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMainHandler(t *testing.T) {
	// Создаем тестовый HTTP-запрос
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()
	handler := middleware.CORS(routers.NewRouter())

	// Вызываем обработчик
	handler.ServeHTTP(rr, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusOK, rr.Code, "Ожидался статус код 200")
}

func TestServerShutdown(t *testing.T) {
	// Создаем тестовый сервер
	server := &http.Server{
		Addr:    ":8081", // Используем другой порт, чтобы не конфликтовать с основным сервером
		Handler: middleware.CORS(routers.NewRouter()),
	}

	// Запускаем сервер в горутине
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Даем серверу время запуститься
	time.Sleep(100 * time.Millisecond)

	// Выполняем graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	assert.NoError(t, err, "Ожидалось успешное завершение работы сервера")
}


func TestCORSHeaders(t *testing.T) {
    req, err := http.NewRequest("OPTIONS", "/", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := middleware.CORS(routers.NewRouter())

    handler.ServeHTTP(rr, req)

    // Обновлённые проверки
    assert.Equal(t, "http://localhost:3000", rr.Header().Get("Access-Control-Allow-Origin"), "Неверный CORS заголовок Origin")
    assert.Equal(t, "Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"), "Неверный CORS заголовок разрешенных заголовков")
}
