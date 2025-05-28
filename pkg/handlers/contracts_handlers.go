package handlers

import (
	db "appContract/pkg/db/repository"
	"appContract/pkg/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func GetAllContracts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method GetAllContracts", http.StatusBadRequest)
		return
	}

	contracts, err := db.DBgetContractAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var contractsResponse []map[string]interface{}
	for _, contract := range contracts {
		var tags []map[string]interface{}
		for _, tag := range contract.Tags {
			tags = append(tags, map[string]interface{}{
				"id_teg":   tag.Id_tags,
				"name_teg": tag.Name_tags,
			})
		}

		contractResponse := map[string]interface{}{
			"contract_id":          contract.Id_contract,
			"name_contract":        contract.Name_contract,
			"surname":              contract.Surname,
			"username":             contract.Username,
			"patronymic":           contract.Patronymic,
			"date_conclusion":      contract.Date_conclusion,
			"date_contract_create": contract.Date_contract_create,
			"date_end":             contract.Date_end,
			"name_type_contract":   contract.Name_type,
			"id_counterparty":      contract.Id_counterparty,
			"name_counterparty":    contract.Name_counterparty,
			"name_status_contract": contract.Name_status_contract,
			"tegs":                 tags,
		}
		contractsResponse = append(contractsResponse, contractResponse)
	}

	data, err := json.Marshal(contractsResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
func GetAllContractsByType(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method GetAllContracts", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	idType := vars["idType"]
	if idType == "" {
		http.Error(w, "idType обязательный параметр", http.StatusBadRequest)
		return
	}
	idTypeInt, err := strconv.Atoi(idType)
	if err != nil {
		http.Error(w, "Недопустимый idType", http.StatusBadRequest)
		return
	}

	contracts, err := db.DBgetContractByType(idTypeInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var contractsResponse []map[string]interface{}
	for _, contract := range contracts {
		var tags []map[string]interface{}
		for _, tag := range contract.Tags {
			tags = append(tags, map[string]interface{}{
				"id_teg":   tag.Id_tags,
				"name_teg": tag.Name_tags,
			})
		}
		contractResponse := map[string]interface{}{
			"contract_id":          contract.Id_contract,
			"name_contract":        contract.Name_contract,
			"surname":              contract.Surname,
			"username":             contract.Username,
			"patronymic":           contract.Patronymic,
			"date_conclusion":      contract.Date_conclusion,
			"date_contract_create": contract.Date_contract_create,
			"date_end":             contract.Date_end,
			"name_type_contract":   contract.Name_type,
			"name_counterparty":    contract.Name_counterparty,
			"name_status_contract": contract.Name_status_contract,
			"tegs":                 tags,
		}
		contractsResponse = append(contractsResponse, contractResponse)
	}
	data, err := json.Marshal(contractsResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func PostAllContractsByDateCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method PostAllContractsByDateCreate", http.StatusBadRequest)
		return
	}
	// Объявляем структуру для входящих данных
	var dateRange struct {
		Date_start string `json:"date_start"`
		Date_end   string `json:"date_end"`
	}

	// Декодируем весь JSON
	err := json.NewDecoder(r.Body).Decode(&dateRange)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Парсим даты
	dateStart, err := time.Parse(time.RFC3339, dateRange.Date_start)
	if err != nil {
		http.Error(w, "Invalid date_start format", http.StatusBadRequest)
		return
	}

	dateEnd, err := time.Parse(time.RFC3339, dateRange.Date_end)
	if err != nil {
		http.Error(w, "Invalid date_end format", http.StatusBadRequest)
		return
	}

	// Создаем объект models.Date
	date := models.Date{
		Date_start: dateStart,
		Date_end:   dateEnd,
	}
	contracts, err := db.DBgetContractsByDateCreate(date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var contractsResponse []map[string]interface{}
	for _, contract := range contracts {
		var tags []map[string]interface{}
		for _, tag := range contract.Tags {
			tags = append(tags, map[string]interface{}{
				"id_teg":   tag.Id_tags,
				"name_teg": tag.Name_tags,
			})
		}
		contractResponse := map[string]interface{}{
			"contract_id":          contract.Id_contract,
			"name_contract":        contract.Name_contract,
			"surname":              contract.Surname,
			"username":             contract.Username,
			"patronymic":           contract.Patronymic,
			"date_conclusion":      contract.Date_conclusion,
			"date_contract_create": contract.Date_contract_create,
			"date_end":             contract.Date_end,
			"name_type_contract":   contract.Name_type,
			"name_counterparty":    contract.Name_counterparty,
			"name_status_contract": contract.Name_status_contract,
			"tegs":                 tags,
		}
		contractsResponse = append(contractsResponse, contractResponse)
	}
	data, err := json.Marshal(contractsResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func GetAllContractsByTegs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method GetAllContractsByTegs", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	idTeg := vars["id_teg_contract"]
	if idTeg == "" {
		http.Error(w, "id_teg_contract обязательный параметр", http.StatusBadRequest)
		return
	}
	idTegInt, err := strconv.Atoi(idTeg)
	if err != nil {
		http.Error(w, "Недопустимый id_teg_contract", http.StatusBadRequest)
		return
	}

	contracts, err := db.DBgetContractByType(idTegInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var contractsResponse []map[string]interface{}
	for _, contract := range contracts {

		var tags []map[string]interface{}
		for _, tag := range contract.Tags {
			tags = append(tags, map[string]interface{}{
				"id_teg":   tag.Id_tags,
				"name_teg": tag.Name_tags,
			})
		}
		contractResponse := map[string]interface{}{
			"contract_id":          contract.Id_contract,
			"name_contract":        contract.Name_contract,
			"surname":              contract.Surname,
			"username":             contract.Username,
			"patronymic":           contract.Patronymic,
			"date_conclusion":      contract.Date_conclusion,
			"date_contract_create": contract.Date_contract_create,
			"date_end":             contract.Date_end,
			"name_type_contract":   contract.Name_type,
			"name_counterparty":    contract.Name_counterparty,
			"name_status_contract": contract.Name_status_contract,
			"tegs":                 tags,
		}
		contractsResponse = append(contractsResponse, contractResponse)
	}
	data, err := json.Marshal(contractsResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func GetAllContractsByStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method GetAllContractsByStatus", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	id_status_contract := vars["id_status_contract"]
	if id_status_contract == "" {
		http.Error(w, "id_status_contract обязательный параметр", http.StatusBadRequest)
		return
	}
	idStatusInt, err := strconv.Atoi(id_status_contract)
	if err != nil {
		http.Error(w, "Недопустимый id_status_contract", http.StatusBadRequest)
		return
	}

	contracts, err := db.DBgetContractByType(idStatusInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var contractsResponse []map[string]interface{}
	for _, contract := range contracts {
		var tags []map[string]interface{}
		for _, tag := range contract.Tags {
			tags = append(tags, map[string]interface{}{
				"id_teg":   tag.Id_tags,
				"name_teg": tag.Name_tags,
			})
		}
		contractResponse := map[string]interface{}{
			"contract_id":          contract.Id_contract,
			"name_contract":        contract.Name_contract,
			"surname":              contract.Surname,
			"username":             contract.Username,
			"patronymic":           contract.Patronymic,
			"date_conclusion":      contract.Date_conclusion,
			"date_contract_create": contract.Date_contract_create,
			"date_end":             contract.Date_end,
			"name_type_contract":   contract.Name_type,
			"name_counterparty":    contract.Name_counterparty,
			"name_status_contract": contract.Name_status_contract,
			"tegs":                 tags,
		}
		contractsResponse = append(contractsResponse, contractResponse)
	}
	data, err := json.Marshal(contractsResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func GetContractID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method GetContract", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	contractId := vars["contractID"]
	if contractId == "" {
		http.Error(w, "Invalid contract_id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(contractId)
	if err != nil {
		http.Error(w, "Invalid contract_id", http.StatusBadRequest)
		return
	}

	contracts, err := db.DBgetContractID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(contracts) == 0 {
		http.Error(w, "Contract not found", http.StatusNotFound)
		return
	}

	contract := contracts[0]

	var tags []map[string]interface{}
	for _, tag := range contract.Tags {
		tags = append(tags, map[string]interface{}{
			"id_teg":   tag.Id_tags,
			"name_teg": tag.Name_tags,
		})
	}

	contractResponse := map[string]interface{}{
		"contract_id":          contract.Id_contract,
		"name_contract":        contract.Name_contract,
		"date_create_contract": contract.Date_contract_create,
		"user_id":              contract.Id_user,
		"username":             contract.Username,
		"surname":              contract.Surname,
		"patronymic":           contract.Patronymic,
		"date_conclusion":      contract.Date_conclusion,
		"date_end":             contract.Date_end,
		"name_type_contract":   contract.Name_type,
		"id_counterparty":      contract.Id_counterparty,
		"name_counterparty":    contract.Name_counterparty,
		"name_status_contract": contract.Name_status_contract,
		"notes":                contract.Notes,
		"conditions":           contract.Conditions,
		"cost":                 contract.Cost,
		"object_contract":      contract.Object_contract,
		"term_payment":         contract.Term_payment,
		"tegs":                 tags,
	}

	data, err := json.Marshal(contractResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func GetUserIDContracts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method GetUserIDContracts", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	userId := vars["userID"]
	if userId == "" {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	contracts, err := db.DBgetContractUserId(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var contractsResponse []map[string]interface{}
	for _, contract := range contracts {

		var tags []map[string]interface{}
		for _, tag := range contract.Tags {
			tags = append(tags, map[string]interface{}{
				"id_teg":   tag.Id_tags,
				"name_teg": tag.Name_tags,
			})
		}

		contractResponse := map[string]interface{}{
			"contract_id":          contract.Id_contract,
			"name_contract":        contract.Name_contract,
			"surname":              contract.Surname,
			"username":             contract.Username,
			"patronymic":           contract.Patronymic,
			"date_conclusion":      contract.Date_conclusion,
			"date_contract_create": contract.Date_contract_create,
			"date_end":             contract.Date_end,
			"name_type_contract":   contract.Name_type,
			"name_counterparty":    contract.Name_counterparty,
			"name_status_contract": contract.Name_status_contract,
			"tegs":                 tags,
		}
		contractsResponse = append(contractsResponse, contractResponse)
	}

	data, err := json.Marshal(contractsResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func PostCreateContract(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method PostCreateContract", http.StatusBadRequest)
		return
	}

	var contract models.Contracts
	err := json.NewDecoder(r.Body).Decode(&contract)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body PostCreateContract", http.StatusBadRequest)
		return
	}

	id, err := db.DBaddContract(contract)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Contract created successfully",
		"id":      id,
	})
}

func PutChangeContract(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	contractId := vars["contractID"]
	id, err := strconv.Atoi(contractId)
	if err != nil {
		http.Error(w, "Invalid contract ID", http.StatusBadRequest)
		return
	}

	var updateData models.Contracts
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	contract := models.Contracts{
		Id_contract:          id,
		Name_contract:        updateData.Name_contract,
		Date_conclusion:      updateData.Date_conclusion,
		Date_end:             updateData.Date_end,
		Id_type:              updateData.Id_type,
		Cost:                 updateData.Cost,
		Object_contract:      updateData.Object_contract,
		Term_payment:         updateData.Term_payment,
		Id_counterparty:      updateData.Id_counterparty,
		Id_status_contract:   updateData.Id_status_contract,
		Notes:                updateData.Notes,
		Conditions:           updateData.Conditions,
	}

	err = db.DBchangeContract(contract)
	if err != nil {
		log.Printf("DB error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Contract updated"})
}

func PutChangeContractUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Id_contract int `json:"id_contract"`
		Id_user     int `json:"id_user"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if request.Id_contract <= 0 {
		http.Error(w, "Contract ID must be positive", http.StatusBadRequest)
		return
	}

	if request.Id_user <= 0 {
		http.Error(w, "User ID must be positive", http.StatusBadRequest)
		return
	}

	if err := db.DBchangeContractUser(request.Id_contract, request.Id_user); err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Contract %d updated with user %d", request.Id_contract, request.Id_user),
	})
}

func DeleteContract(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method - expected DELETE", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	contractID, err := strconv.Atoi(vars["contractID"])
	if err != nil {
		http.Error(w, "Invalid contract ID format", http.StatusBadRequest)
		return
	}

	err = db.DBdeleteContract(contractID)
	if err != nil {

		log.Printf("Error deleting contract %d: %v", contractID, err)

		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Contract not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete contract and related data", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"success":     true,
		"message":     fmt.Sprintf("Contract %d and all related data deleted successfully", contractID),
		"contract_id": contractID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
