package routers

import (
	"appContract/pkg/handlers"

	"github.com/gorilla/mux"
)

func ContractsAndStageRoutes(router *mux.Router) {
    // Контракты
    router.HandleFunc("/api/contractsAndStage", handlers.GetContractsAndStages).Methods("GET")
	//router.HandleFunc("/api/contractsAndStage", handlers.GetContractsandStags).Methods("GET")
  
    
}