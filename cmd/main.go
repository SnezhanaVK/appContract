package main

import (
	"appContract/pkg/db"
	"appContract/pkg/middleware"
	"appContract/pkg/routers"
	"context"
	"os"
	"os/signal"
	"syscall"

	"log"
	"net/http"
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
	if err:=server.Shutdown(context.Background()); err!=nil{
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server exited properly")
}


	
