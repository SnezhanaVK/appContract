// contracts.go
package handlers

import (
	"appContract/pkg/models"
	"encoding/json"
	"net/http"
	"strconv"

	db "appContract/pkg/db/repository"
)

// Contracts
func GetAllContracts(w http.ResponseWriter, r *http.Request) {
    if r.Method!=http.MethodGet{
        http.Error(w,"Invalid request method GetAllContracts",http.StatusBadRequest)
        return
    }
    contracts, err := db.DBgetContractAll()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(contracts)
}

func GetUserIDContracts(w http.ResponseWriter, r *http.Request) {

    if r.Method!=http.MethodGet{
        http.Error(w,"Invalid request method GetUserIDContracts",http.StatusBadRequest)
        return
    }
    userId, err:= strconv.Atoi(r.URL.Query().Get("user_id"))
    if err != nil {
        http.Error(w, "Invalid user_id", http.StatusBadRequest)
        return
    }

    contracts, err := db.DBgetContractUserId(userId)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(contracts)
}

func GetContractID(w http.ResponseWriter, r *http.Request) {
    if r.Method!=http.MethodGet{
        http.Error(w,"Invalid request method GetContract",http.StatusBadRequest)
        return
    }
    constantId, err:= strconv.Atoi(r.URL.Query().Get("contract_id"))
    if err != nil {
        http.Error(w, "Invalid contract_id", http.StatusBadRequest)
        return
    }
    contract, err := db.DBgetContractID(constantId)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(contract)
}

func PostCreateContract(w http.ResponseWriter, r *http.Request) {
  if r.Method!=http.MethodPost{
        http.Error(w,"Invalid request method PostCreateContract",http.StatusBadRequest)
        return
    }
    var contract models.Contracts
    err:=json.NewDecoder(r.Body).Decode(&contract)
    if err!=nil{
        http.Error(w,"Invalid request body PostCreateContract",http.StatusBadRequest)
        return
    }
    err = db.DBaddContract(contract)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Contract created successfully"})
}

func PutChangeContract(w http.ResponseWriter, r *http.Request) {
    if r.Method!=http.MethodPut{
        http.Error(w,"Invalid request method PutChangeContract",http.StatusBadRequest)
        return
    }
    var contract models.Contracts
    err:=json.NewDecoder(r.Body).Decode(&contract)
    if err!=nil{
        http.Error(w,"Invalid request body PutChangeContract",http.StatusBadRequest)
        return
    }
    err = db.DBchangeContract(contract)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Contract updated successfully"})
}


func PutChangeContractUser(w http.ResponseWriter, r *http.Request) {
    if r.Method!=http.MethodPut{
        http.Error(w,"Invalid request method UpdateContractUser",http.StatusBadRequest)
        return
    }
    userId, err:= strconv.Atoi(r.URL.Query().Get("user_id"))
    contract_id, err:= strconv.Atoi(r.URL.Query().Get("contract_id"))
    if err != nil {
        http.Error(w, "Invalid user_id", http.StatusBadRequest)
        return
    }
    var contract models.Contracts
    err=json.NewDecoder(r.Body).Decode(&contract)
    if err!=nil{
        http.Error(w,"Invalid request body UpdateContractUser",http.StatusBadRequest)
        return
    }
    err = db.DBchangeContractUser(userId,contract_id)
    if err != nil { 
        http.Error(w, err.Error(), http.StatusInternalServerError)        
        return
        
    }
}

func DeleteContract(w http.ResponseWriter, r *http.Request) {
    if r.Method!=http.MethodDelete{
        http.Error(w,"Invalid request method DeleteContract",http.StatusBadRequest)
        return
    }
    constantId, err:= strconv.Atoi(r.URL.Query().Get("contract_id"))
    if err != nil {
        http.Error(w, "Invalid contract_id", http.StatusBadRequest)
        return
    }
    err =db.DBdeleteContract(constantId)
    if err !=nil{
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Contract deleted successfully"})
}
