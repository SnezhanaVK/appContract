package main

import (
	"appContract/pkg/db"
	"appContract/pkg/routers"
	"log"
	"net/http"
)


func main() {

    db.SetupDatabase()

    router := routers.NewRouter()
    log.Println("Сервер запущен на порту :8080")
    log.Println("http://localhost:8080")

    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatal(err)
    }
    
}