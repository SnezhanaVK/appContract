package routers

import (
	"appContract/pkg/handlers"

	"github.com/gorilla/mux"
)

func UserRoutes(router *mux.Router) {
    // Пользователи
    router.HandleFunc("/api/users", handlers.GetAllUsers).Methods("GET")
    router.HandleFunc("/api/users/{userID}", handlers.GetUserID).Methods("GET")
    router.HandleFunc("/api/users/create", handlers.PostCreateUser).Methods("POST")
    router.HandleFunc("/api/users/{id}", handlers.PutUpdateUser).Methods("PUT")
    router.HandleFunc("/api/users/{id}", handlers.DeleteUser).Methods("DELETE")
}
