package routers

import (
	"appContract/pkg/handlers"

	"github.com/gorilla/mux"
)

func NotificationRoutes(router *mux.Router) {
	router.HandleFunc("/api/users/{userID}/notifications", handlers.GetSettings).Methods("GET")
	router.HandleFunc("/api/users/{userID}/notifications", handlers.PutUpdateUserSettings).Methods("PUT")
	router.HandleFunc("/api/users/{userID}/notifications", handlers.DeleteSettings).Methods("DELETE")
}
