package handlers

import (
	db "appContract/pkg/db/repository"
	"encoding/json"
	"net/http"
)

func GetContractsandStags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method GetContractsandStags", http.StatusBadRequest)
		return
	}
	contracts, err := db.DBgetContractAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var contractsAndStagesResponse []map[string]interface{}
	for _, contract := range contracts {
		contractResponse := map[string]interface{}{
			"contract_id":          contract.Id_contract,
			"name_contract":        contract.Name_contract,
			"date_create_contract": contract.Date_contract_create,
			"user_id":              contract.Id_user,
			"date_conclusion":      contract.Date_conclusion,
			"date_start":           contract.Date_contract_create,
			"date_end":             contract.Date_end,
			"id_type":              contract.Id_type,
			"name_type_contract":   contract.Name_type,
			"id_counterparty":      contract.Id_counterparty,
			"name_counterparty":    contract.Name_counterparty,
			"id_status_contract":   contract.Id_status_contract,
			"name_status_contract": contract.Name_status_contract,
			"id_teg":               contract.Id_teg_contract,
			"name_teg":             contract.Tegs_contract,
		}
		contractsAndStagesResponse = append(contractsAndStagesResponse, contractResponse)
		stages, err := db.DBgetStageAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, stage := range stages {

			stageResponse := map[string]interface{}{
				"id_stage":             stage.Id_stage,
				"name_stage":           stage.Name_stage,
				"id_user":              stage.Id_user,
				"surname":              stage.Surname,
				"username":             stage.Username,
				"patronymic":           stage.Patronymic,
				"phone":                stage.Phone,
				"email":                stage.Email,
				"description":          stage.Description,
				"status_stage":         stage.Id_status_stage,
				"date_change_status":   stage.Date_change_status,
				"name_status_stage":    stage.Name_status_stage,
				"date_create_start":    stage.Date_create_start,
				"date_create_end":      stage.Date_create_end,
				"id_contract":          stage.Id_contract,
				"name_contract":        stage.Name_contract,
				"date_create_contract": stage.Data_contract_create,
			}
			contractsAndStagesResponse = append(contractsAndStagesResponse, stageResponse)
		}
	}
	data, err := json.Marshal(contractsAndStagesResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
