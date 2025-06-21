// main.go
package main

import (
	"appContract/pkg/db"
	"appContract/pkg/middleware"
	"appContract/pkg/routers"
	"appContract/pkg/service"
	"appContract/pkg/utils"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	err := db.SetupDatabase()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.CloseDB()
	db.ConnectDB()

	emailSender := utils.NewDefaultEmailSender()
	service.InitEmailSender(emailSender)

	router := routers.NewRouter()
	handlers := middleware.CORS(router)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handlers,
		    ReadTimeout:  15 * time.Second,
    WriteTimeout: 30 * time.Second,
    IdleTimeout:  60 * time.Second,
    MaxHeaderBytes: 1 << 20, // 1MB

	}

	c := cron.New()
	_, err = c.AddFunc("15 17 * * *", func() {
		log.Println("Запуск обработки уведомлений...")
		if err := service.ProcessDailyNotifications(); err != nil {
			log.Printf("Ошибка обработки уведомлений: %v", err)
		} else {
			log.Println("Обработка уведомлений завершена успешно")
		}
	})
	if err != nil {
		log.Fatalf("Ошибка настройки расписания: %v", err)
	}
	c.Start()

	go func() {
		log.Println("http://localhost:8080")
		log.Println("Monitoring: " + "http://localhost:8080/metrics")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server is shutting down...")
	c.Stop()

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server exited properly")
}
