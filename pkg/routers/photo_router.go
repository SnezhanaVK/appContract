package routers

import (
	"appContract/pkg/handlers"

	"github.com/gorilla/mux"
)

func PhotoRouters(router *mux.Router) {
	router.HandleFunc("/api/photo/{user_id}", handlers.PostAddPhoto).Methods("POST")
	router.HandleFunc("/api/photo/{user_id}", handlers.GetPhoto).Methods("GET")
	router.HandleFunc("/api/photo/{user_id}", handlers.DeletePhoto).Methods("DELETE")
	router.HandleFunc("/api/photo/{user_id}", handlers.PutChangePhoto).Methods("PUT")
}
