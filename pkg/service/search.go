package service

import (
	"appContract/pkg/db"
	"appContract/pkg/models"
	"context"

	"log"
)

func SearchContract(nameContract string, nameStage string, nameTeg string) []models.Contracts {

	conn := db.GetDB()
	if conn == nil {
		return nil
	}

	rows, err := conn.Query(context.Background(), "SELECT * FROM contracts WHERE name_contract LIKE $1 OR name_stage IN (SELECT name_stage FROM stages WHERE name_stage LIKE $2) OR name_teg IN (SELECT name_teg FROM tegs WHERE name_teg LIKE $3)", "%"+nameContract+"%", "%"+nameStage+"%", "%"+nameTeg+"%")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var contracts []models.Contracts
	for rows.Next() {
		var contract models.Contracts
		err := rows.Scan(&contract.Id_contract, &contract.Name_contract, &contract.Date_contract_create, &contract.Id_user, &contract.Date_conclusion, &contract.Date_end, &contract.Id_type, &contract.Cost, &contract.Object_contract, &contract.Term_payment, &contract.Id_counterparty, &contract.Id_status_contract, &contract.Notes, &contract.Conditions)
		if err != nil {
			log.Fatal(err)
		}
		contracts = append(contracts, contract)
	}
	return contracts
}
