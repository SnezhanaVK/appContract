// contracts.go
package handlers

import (
	"appContract/pkg/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	db "appContract/pkg/db/repository"

	"github.com/gorilla/mux"
)

// Contracts
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
        // Создаем массив для тегов
        var tegs []map[string]interface{}
        for _, teg := range contract.Tegs {
            tegs = append(tegs, map[string]interface{}{
               /// "id_teg":   teg.Id_tegs,
                "name_teg": teg.Name_tegs,
            })
        }

        contractResponse := map[string]interface{}{
            "contract_id":          contract.Id_contract,
            "name_contract":       contract.Name_contract,
           
            "surname":             contract.Surname,
            "username":            contract.Username,
            "patronymic":          contract.Patronymic,
            "date_conclusion":     contract.Date_conclusion,
            "date_contract_create": contract.Date_contract_create,
            "date_end":           contract.Date_end,
            "name_type_contract":  contract.Name_type,
            "name_counterparty":   contract.Name_counterparty,
            "name_status_contract": contract.Name_status_contract,
            "tegs":               tegs, // Включаем массив тегов
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
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    vars := mux.Vars(r)
    idType := vars["idType"]
    if idType == "" {
        http.Error(w, "idType is required", http.StatusBadRequest)
        return
    }

    idTypeInt, err := strconv.Atoi(idType)
    if err != nil {
        http.Error(w, "Invalid idType", http.StatusBadRequest)
        return
    }

    contracts, err := db.DBgetContractByType(idTypeInt)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var contractsResponse []map[string]interface{}
    for _, contract := range contracts {
        // Создаем массив для тегов
        var tags []map[string]interface{}
        for _, tag := range contract.Tegs {
            tags = append(tags, map[string]interface{}{
                "id":   tag.Id_tegs,
                "name": tag.Name_tegs,
            })
        }

        contractResponse := map[string]interface{}{
            "id_contract":          contract.Id_contract,
            "name_contract":        contract.Name_contract,
            "date_contract_create": contract.Date_contract_create,
            "date_conclusion":      contract.Date_conclusion,
            "date_end":            contract.Date_end,
            "user": map[string]interface{}{
                "id":         contract.Id_user,
                "surname":    contract.Surname,
                "name":       contract.Username,
                "patronymic": contract.Patronymic,
            },
            "type": map[string]interface{}{
                "id":   contract.Id_type,
                "name": contract.Name_type,
            },
            "counterparty": map[string]interface{}{
                "id":   contract.Id_counterparty,
                "name": contract.Name_counterparty,
            },
            "status": map[string]interface{}{
                "id":   contract.Id_status_contract,
                "name": contract.Name_status_contract,
            },
            "tags": tags,
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

func PostAllContractsByDateCreate(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Получаем даты из тела запроса
    var dateRange struct {
        Date_start time.Time `json:"date_start"`
        Date_end   time.Time `json:"date_end"`
    }
    
    err := json.NewDecoder(r.Body).Decode(&dateRange)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Создаем структуру Date для запроса к БД
    date := models.Date{
        Date_start: dateRange.Date_start,
        Date_end:   dateRange.Date_end,
    }

    contracts, err := db.DBgetContractsByDateCreate(date)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var contractsResponse []map[string]interface{}
    for _, contract := range contracts {
        // Подготавливаем массив тегов
        var tags []map[string]interface{}
        for _, tag := range contract.Tegs {
            tags = append(tags, map[string]interface{}{
                "id_tegs":   tag.Id_tegs,
                "name_tegs": tag.Name_tegs,
            })
        }

        contractResponse := map[string]interface{}{
            "contract_id":          contract.Id_contract,
            "name_contract":        contract.Name_contract,
            "date_create_contract": contract.Date_contract_create,
            "user": map[string]interface{}{
                "id":         contract.Id_user,
                "surname":    contract.Surname,
                "name":       contract.Username,
                "patronymic": contract.Patronymic,
            },
            "date_conclusion":     contract.Date_conclusion,
            "date_end":           contract.Date_end,
            "type": map[string]interface{}{
                "id":   contract.Id_type,
                "name": contract.Name_type,
            },
            "counterparty": map[string]interface{}{
                "id":   contract.Id_counterparty,
                "name": contract.Name_counterparty,
            },
            "status": map[string]interface{}{
                "id":   contract.Id_status_contract,
                "name": contract.Name_status_contract,
            },
            "tags": tags,
        }
        contractsResponse = append(contractsResponse, contractResponse)
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(contractsResponse); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
}
func GetAllContractsByTegs(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Получаем контракты с тегами из БД
    contracts, err := db.DBgetContractsByTegs()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var contractsResponse []map[string]interface{}
    for _, contract := range contracts {
        // Подготавливаем массив тегов
        var tags []map[string]interface{}
        for _, tag := range contract.Tegs {
            tags = append(tags, map[string]interface{}{
                "id_teg":   tag.Id_tegs,
                "name_teg": tag.Name_tegs,
            })
        }

        contractResponse := map[string]interface{}{
            "contract_id":          contract.Id_contract,
            "name_contract":       contract.Name_contract,
            "date_create_contract": contract.Date_contract_create,
            "user": map[string]interface{}{
                "id":         contract.Id_user,
                "surname":    contract.Surname,
                "name":       contract.Username,
                "patronymic": contract.Patronymic,
            },
            "date_conclusion":     contract.Date_conclusion,
            "date_end":           contract.Date_end,
            "type": map[string]interface{}{
                "id":   contract.Id_type,
                "name": contract.Name_type,
            },
            "counterparty": map[string]interface{}{
                "id":   contract.Id_counterparty,
                "name": contract.Name_counterparty,
            },
            "status": map[string]interface{}{
                "id":   contract.Id_status_contract,
                "name": contract.Name_status_contract,
            },
            "tags": tags,
        }
        contractsResponse = append(contractsResponse, contractResponse)
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(contractsResponse); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
}

func GetAllContractsByStatus(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    vars := mux.Vars(r)
    idStatusStr := vars["id_status_contract"]
    if idStatusStr == "" {
        http.Error(w, "Status ID is required", http.StatusBadRequest)
        return
    }

    idStatus, err := strconv.Atoi(idStatusStr)
    if err != nil {
        http.Error(w, "Invalid status ID", http.StatusBadRequest)
        return
    }

    contracts, err := db.DBgetContractsByStatus(idStatus)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var response []map[string]interface{}
    for _, contract := range contracts {
        // Подготавливаем теги
        var tags []map[string]interface{}
        for _, tag := range contract.Tegs {
            tags = append(tags, map[string]interface{}{
                "id":   tag.Id_tegs,
                "name": tag.Name_tegs,
            })
        }

        // Формируем структуру ответа
        contractData := map[string]interface{}{
            "id": contract.Id_contract,
            "name": contract.Name_contract,
            "dates": map[string]interface{}{
                "conclusion": contract.Date_conclusion,
                "create":     contract.Date_contract_create,
                "end":       contract.Date_end,
            },
            "user": map[string]interface{}{
                "id":         contract.Id_user,
                "surname":    contract.Surname,
                "name":       contract.Username,
                "patronymic": contract.Patronymic,
            },
            "type": map[string]interface{}{
                "id":   contract.Id_type,
                "name": contract.Name_type,
            },
            "counterparty": map[string]interface{}{
                "id":   contract.Id_counterparty,
                "name": contract.Name_counterparty,
            },
            "status": map[string]interface{}{
                "id":   contract.Id_status_contract,
                "name": contract.Name_status_contract,
            },
            "tags": tags,
        }
        response = append(response, contractData)
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
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

    // Проверяем, что контракт найден
    if len(contracts) == 0 {
        http.Error(w, "Contract not found", http.StatusNotFound)
        return
    }

    // Берем первый контракт (должен быть только один)
    contract := contracts[0]

    // Создаем массив для тегов
    var tegs []map[string]interface{}
    for _, teg := range contract.Tegs {
        tegs = append(tegs, map[string]interface{}{
            "id_teg":   teg.Id_tegs,
            "name_teg": teg.Name_tegs,
        })
    }

    contractResponse := map[string]interface{}{
        "contract_id":          contract.Id_contract,
        "name_contract":        contract.Name_contract,
        "date_create_contract": contract.Date_contract_create,
        "user_id":             contract.Id_user,
        "username":             contract.Username,
        "surname":              contract.Surname,
        "patronymic":           contract.Patronymic,
        "date_conclusion":     contract.Date_conclusion,
        "date_end":            contract.Date_end,
        "name_type_contract":  contract.Name_type,
        "name_counterparty":   contract.Name_counterparty,
        "name_status_contract": contract.Name_status_contract,
        "tegs":                tegs, // Добавляем массив тегов
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

func GetUserContracts(w http.ResponseWriter, r *http.Request) {
    // Проверка метода запроса
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Извлечение и валидация ID пользователя
    vars := mux.Vars(r)
    userID := vars["userID"]
    if userID == "" {
        http.Error(w, "User ID is required", http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(userID)
    if err != nil {
        http.Error(w, "Invalid user ID format", http.StatusBadRequest)
        return
    }

    // Получение контрактов пользователя
    contracts, err := db.DBgetContractUserId(id)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to get user contracts: %v", err), 
            http.StatusInternalServerError)
        return
    }

    // Формирование ответа
    var response []map[string]interface{}
    for _, contract := range contracts {
        // Подготовка тегов
        var tags []map[string]interface{}
        for _, tag := range contract.Tegs {
            tags = append(tags, map[string]interface{}{
                "id":   tag.Id_tegs,
                "name": tag.Name_tegs,
            })
        }

        // Структура контракта
        contractData := map[string]interface{}{
            "id":   contract.Id_contract,
            "name": contract.Name_contract,
            "dates": map[string]interface{}{
                "create":      contract.Date_contract_create,
                "conclusion":  contract.Date_conclusion,
                "end":         contract.Date_end,
            },
            "type": map[string]interface{}{
                "id":   contract.Id_type,
                "name": contract.Name_type,
            },
            "counterparty": map[string]interface{}{
                "id":   contract.Id_counterparty,
                "name": contract.Name_counterparty,
            },
            "status": map[string]interface{}{
                "id":   contract.Id_status_contract,
                "name": contract.Name_status_contract,
            },
            "tags": tags,
        }
        response = append(response, contractData)
    }

    // Отправка ответа
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, "Failed to encode response", 
            http.StatusInternalServerError)
    }
}



func PostCreateContract(w http.ResponseWriter, r *http.Request) {
  if r.Method!=http.MethodPost{
        http.Error(w,"Invalid request method PostCreateContract",http.StatusBadRequest)
        return
    }
    var contract models.Contracts
    err:=json.NewDecoder(r.Body).Decode(&contract)
    if err!=nil{
        log.Println(err)
        http.Error(w,"Invalid request body PostCreateContract",http.StatusBadRequest)
        return
    }
    err = db.DBaddContract(contract)
    if err != nil {
        log.Println(err)
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
    if r.Method != http.MethodPut {
        http.Error(w, "Invalid request method UpdateContractUser", http.StatusBadRequest)
        return
    }

    var data map[string]int
    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    contract_id, ok := data["id_contract"]
    if !ok {
        http.Error(w, "Missing id_contract in request body", http.StatusBadRequest)
        return
    }

    userId, ok := data["id_user"]
    if !ok {
        http.Error(w, "Missing id_user in request body", http.StatusBadRequest)
        return
    }

    err = db.DBchangeContractUser(contract_id, userId)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Contract user updated successfully"})
}

func DeleteContract(w http.ResponseWriter, r *http.Request) {
    if r.Method!=http.MethodDelete{
        http.Error(w,"Invalid request method DeleteContract",http.StatusBadRequest)
        return
    }
    vars:=mux.Vars(r)
    contractId, err := strconv.Atoi(vars["contractID"])
    if err != nil {
        http.Error(w, "Invalid contract_id", http.StatusBadRequest)
        return
    }
    err =db.DBdeleteContract(contractId)
    if err !=nil{
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Contract deleted successfully"})
}
