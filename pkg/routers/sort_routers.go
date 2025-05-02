package routers

import (
	"appContract/pkg/handlers"

	"github.com/gorilla/mux"
)

func SortRouters(router *mux.Router) {
	router.HandleFunc("/api/sort/statusContract", handlers.GetStatusContract).Methods("GET")
	router.HandleFunc("/api/sort/statusStage", handlers.GetStatusStage).Methods("GET")
	router.HandleFunc("/api/sort/tags", handlers.GetTags).Methods("Get")
	router.HandleFunc("/api/sort/types", handlers.GetType).Methods("GET")
	router.HandleFunc("/api/contracts/{contractId}/tags", handlers.AddTagToContractHandler).Methods("POST")
	router.HandleFunc("/api/contracts/{contractId}/tags/{tagId}", handlers.RemoveTagFromContractHandler).Methods("DELETE")
	router.HandleFunc("/api/contracts/{contractId}/tags", handlers.GetContractTagsHandler).Methods("GET")

}
