package routers

import (
	"appContract/pkg/handlers"

	"github.com/gorilla/mux"
)

func StageRoutes(router *mux.Router) {
    // Этапы
    router.HandleFunc("/api/stages", handlers.GetAllStages).Methods("GET")
    router.HandleFunc("/api/stages/{userID}", handlers.GetUserStages).Methods("GET")
    router.HandleFunc("/api/stages/{stageID}", handlers.GetStage).Methods("GET")
    router.HandleFunc("/api/stages/{id}/files", handlers.GetStageFiles).Methods("GET")
    router.HandleFunc("/api/stages/{id}/files", handlers.AddFileToStage).Methods("POST")
    router.HandleFunc("/api/stages/create", handlers.CreateStage).Methods("POST")
    router.HandleFunc("/api/stages/{id}/files", handlers.DeleteStageFiles).Methods("DELETE")
    router.HandleFunc("/api/stages/{id}", handlers.DeleteStage).Methods("DELETE")
    router.HandleFunc("/api/stages/{id}/status", handlers.UpdateStageStatus).Methods("PUT")
    router.HandleFunc("/api/stages/{id}/status/{statusID}", handlers.GetStageStatus).Methods("GET")
    router.HandleFunc("/api/stages/{id}/settings", handlers.UpdateNotificationSettings).Methods("PUT")
    router.HandleFunc("/api/stages/{id}/settings/{settingID}", handlers.GetNotificationSettings).Methods("PUT")
    router.HandleFunc("/api/stages/comment", handlers.CreateComment).Methods("POST")
    router.HandleFunc("/api/stages/comment/{id}", handlers.GetComment).Methods("GET")
    router.HandleFunc("/api/stages/{id}/comment", handlers.GetStageComments).Methods("GET")
    router.HandleFunc("/api/stages/informationDate/{stageID}", handlers.GetNotificationInfo).Methods("GET")
}
