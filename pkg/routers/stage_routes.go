package routers

import (
	"appContract/pkg/handlers"

	"github.com/gorilla/mux"
)

func StageRoutes(router *mux.Router) {
	router.HandleFunc("/api/stages", handlers.GetAllStages).Methods("GET")
	router.HandleFunc("/api/stages/userID/{userID}", handlers.GetUserStages).Methods("GET")
	router.HandleFunc("/api/stages/{stageID}", handlers.GetStage).Methods("GET")
	router.HandleFunc("/api/stages/contractId/{contractID}", handlers.GetStagesByIdContract).Methods("GET")
	router.HandleFunc("/api/stages/{stageID}/files/{fileID}", handlers.GetStageFilesID).Methods("GET")
	router.HandleFunc("/api/stages/{stageID}/files", handlers.GetStageFiles).Methods("GET")
	router.HandleFunc("/api/stages/status/{statusID}", handlers.GetStageStatus).Methods("GET")
	router.HandleFunc("/api/stages/{stageID}/comment", handlers.GetComments).Methods("GET")
	router.HandleFunc("/api/stages/{stage_id}/files", handlers.PostFileToStage).Methods("POST")
	router.HandleFunc("/api/stages/create", handlers.PostCreateStage).Methods("POST")
	router.HandleFunc("/api/stages/{stageID}/status/{idStatusStage}/comment", handlers.PostAddComment).Methods("POST")
	router.HandleFunc("/api/stages/{StageID}/status", handlers.PutStageStatus).Methods("PUT")
	router.HandleFunc("/api/stages/files/{id_file}", handlers.DeleteStageFiles).Methods("DELETE")
	router.HandleFunc("/api/stages/{stageID}", handlers.DeleteStage).Methods("DELETE")
	router.HandleFunc("/api/stages/comment/{idComment}", handlers.DeleteComment).Methods("DELETE")

}
