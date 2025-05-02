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
    router.HandleFunc("/api/users/rolesUser/{userID}", handlers.GetUserRoles).Methods("GET")
    router.HandleFunc("/api/users/addRoleAdmin/{userID}", handlers.PostAddRoleAdmin).Methods("POST")
    router.HandleFunc("/api/users/addRoleManager/{userID}", handlers.PostAddRoleManager).Methods("POST")
    router.HandleFunc("/api/users/deleteRoleUser/{userID}", handlers.DeleteRemoveRoleAdmin).Methods("Delete")
    router.HandleFunc("/api/users/deleteRoleManager/{userID}", handlers.DeleteRemoveRoleManager).Methods("Delete")
    router.HandleFunc("/api/users/{id}", handlers.PutUpdateUser).Methods("PUT")
    router.HandleFunc("/api/users/{id}", handlers.DeleteUser).Methods("DELETE")
}
