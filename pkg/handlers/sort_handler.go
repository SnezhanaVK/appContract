package handlers

import (
	db "appContract/pkg/db/repository"
	"encoding/json"
	"net/http"
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
            "id_status_stage":  status.Id_status_stage,
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