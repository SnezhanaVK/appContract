package routers

import (
	"appContract/pkg/handlers"

	"github.com/gorilla/mux"
)

func ContractRoutes(router *mux.Router) {
    // Контракты
    router.HandleFunc("/api/contracts", handlers.GetAllContracts).Methods("GET")
    router.HandleFunc("/api/contracts/user/{userID}", handlers.GetUserIDContracts).Methods("GET")
    router.HandleFunc("/api/contracts/{contractID}", handlers.GetContractID).Methods("GET")
    router.HandleFunc("/api/contracts/create", handlers.PostCreateContract).Methods("POST")
    router.HandleFunc("/api/contracts/{id}", handlers.PutChangeContract).Methods("PUT")
    router.HandleFunc("/api/contracts/userchange", handlers.PutChangeContractUser).Methods("PUT")
    router.HandleFunc("/api/contracts/{contractID}", handlers.DeleteContract).Methods("DELETE")
}
