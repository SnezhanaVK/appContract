package main

import (
	"appContract/pkg/db"
	dbrepo "appContract/pkg/db/repository"
	"appContract/pkg/middleware"
	"appContract/pkg/routers"
	service "appContract/pkg/services"
	"appContract/pkg/utils"
	"context"
	"os"
	"os/signal"
	"syscall"

	"log"
	"net/http"

	"github.com/robfig/cron/v3"
)

func main() {
	err:=db.SetupDatabase()
	if err!=nil{
		log.Fatal("Error connecting to datebase :%v", err)
	}
	defer db.CloseDB()
	db.ConnectDB()

	router := routers.NewRouter()
	handlers:=middleware.CORS(router)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handlers,
	}
	emailSender := utils.NewEmailSender(
		"snezhanakydryavtseva@gmail.com",    // Ваш email
		"ixjl oetf pqop bgos",        // Пароль приложения (для Gmail)
		"smtp.gmail.com",           // SMTP хост
		"587",                      // SMTP порт
	)
	
	notificationRepo := dbrepo.NewNotificationRepository(db.GetDBConection()) // Теперь должно работать
	notificationService := service.NewNotificationService(notificationRepo, emailSender)
	// Настройка cron-расписания
	c := cron.New()
	
	// Запуск задачи ежедневно в 00:00
	_, err = c.AddFunc("0 0 * * *", func() {
		log.Println("Запуск обработки уведомлений...")
		if err := notificationService.ProcessDailyNotifications(); err != nil {
			log.Printf("Ошибка обработки уведомлений: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Ошибка настройки расписания: %v", err)
	}
	
	c.Start()
	
	go func() {
		log.Println("Сервер запущен на порту :8080")
		log.Println("http://localhost:8080")
		if err:=server.ListenAndServe(); 
		err!=nil&&err!=http.ErrServerClosed{
			log.Fatalf("Server error: %v", err)	
		}
	}()
	quit:=make(chan os.Signal,1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server is shutting down...")
	c.Stop()
	if err:=server.Shutdown(context.Background()); err!=nil{
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server exited properly")
}


	
