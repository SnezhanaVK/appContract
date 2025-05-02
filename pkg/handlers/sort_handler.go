package handlers

import (
	db "appContract/pkg/db/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	tags, err := db.DBGetTags()
	if err != nil {
		http.Error(w, "Error fetching data from database", http.StatusInternalServerError)
		return
	}

	var tagsResponse []map[string]interface{}
	for _, tag := range tags {
		tagResponse := map[string]interface{}{
			"Id_teg_contract": tag.Id_teg_contract,
			"Tegs_contract":   tag.Tegs_contract,
		}
		tagsResponse = append(tagsResponse, tagResponse)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tagsResponse); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
	}
}
func GetStatusContract(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	statuses, err := db.DBGetStatusContract()
	if err != nil {
		http.Error(w, "Error fetching data from database", http.StatusInternalServerError)
		return
	}

	var response []map[string]interface{}
	for _, status := range statuses {
		statusResponse := map[string]interface{}{
			"id_status_contract":   status.Id_status_contract,
			"name_status_contract": status.Name_status_contract,
		}
		response = append(response, statusResponse)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
	}
}

func GetStatusStage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed) // Исправлен статус код
		return
	}

	statuses, err := db.DBGetStatusStage()
	if err != nil {
		http.Error(w, "Error fetching data from database", http.StatusInternalServerError)
		return
	}

	// Создаем массив для ответа с нужными полями
	var response []map[string]interface{}
	for _, status := range statuses {
		statusResponse := map[string]interface{}{
			"id_status_stage":   status.Id_status_stage,
			"name_status_stage": status.Name_status_stage,
		}
		response = append(response, statusResponse)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
	}
}

func GetType(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed) // Исправлен статус код
		return
	}

	types, err := db.DBGetTypeContract()
	if err != nil {
		http.Error(w, "Error fetching data from database", http.StatusInternalServerError) // Исправлена опечатка
		return
	}

	// Создаем массив с нужными полями
	var response []map[string]interface{}
	for _, t := range types {
		typeResponse := map[string]interface{}{
			"id_type_contract":   t.Id_type,
			"name_type_contract": t.Name_type,
		}
		response = append(response, typeResponse)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
	}
}

// Обработчики
func AddTagToContractHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contractID, err := strconv.Atoi(vars["contractId"])
	if err != nil {
		http.Error(w, "Invalid contract ID", http.StatusBadRequest)
		return
	}

	var request struct {
		TagID int `json:"tagId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := db.AddTagToContract(contractID, request.TagID); err != nil {
		switch err.Error() {
		case "tag does not exist":
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{
		Message: "Тег успешно добавлен к контракту",
	})
}

func RemoveTagFromContractHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contractID, err := strconv.Atoi(vars["contractId"])
	if err != nil {
		http.Error(w, "Invalid contract ID", http.StatusBadRequest)
		return
	}

	tagID, err := strconv.Atoi(vars["tagId"])
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	if err := db.RemoveTagFromContract(contractID, tagID); err != nil {
		switch err.Error() {
		case "tag association not found":
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		Success: true,
		Message: "Tag successfully removed from contract",
	})
}

func GetContractTagsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contractID, err := strconv.Atoi(vars["contractId"])
	if err != nil {
		http.Error(w, "Invalid contract ID", http.StatusBadRequest)
		return
	}

	tags, err := db.GetContractIDTags(contractID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}
