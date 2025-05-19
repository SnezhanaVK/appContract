package routers

import (
	"appContract/pkg/handlers"

	"github.com/gorilla/mux"
)

func ContractsAndStageRoutes(router *mux.Router) {
	router.HandleFunc("/api/contractsAndStage", handlers.GetContractsandStags).Methods("GET")
}
