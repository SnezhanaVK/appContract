package main

import (
	"appContract/pkg/db"
	"appContract/pkg/middleware"
	"appContract/pkg/routers"
	"log"
	"net/http"
)

func main() {
	db.SetupDatabase()

	router := routers.NewRouter()
	
	// Оберните роутер в CORS middleware
	handler := middleware.CORS(router)
	
	log.Println("Сервер запущен на порту :8080")
	log.Println("http://localhost:8080")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}