package db

import (
	db "appContract/pkg/db"
	"appContract/pkg/models"
	"context"
	"errors"
	"log"
)
func DBGetTags() ([]models.Contracts, error) {
	conn := db.GetDB()
	if conn == nil {
		return nil, errors.New("database connection is nil")
	}
	
	rows, err := conn.Query(context.Background(), "SELECT id_teg, name_teg FROM tegs")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tags []models.Contracts
	for rows.Next() {
		var tag models.Contracts
		err := rows.Scan(&tag.Id_teg_contract,
						 &tag.Tegs_contract)
		if err != nil {
			log.Fatal(err)
		}
		tags = append(tags, tag)
	}
	

	return tags, nil
}
func DBGetStatusContract() ([]models.Contracts, error) {
	conn := db.GetDB()
	if conn == nil {
		return nil, errors.New("database connection is nil")
	}
	
	rows, err := conn.Query(context.Background(), "SELECT id_status_contract, name_status_contract FROM status_contracts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var statuses []models.Contracts
	for rows.Next() {
		var  status models.Contracts
		err := rows.Scan(&status.Id_status_contract, &status.Name_status_contract)
		if err != nil {
			log.Fatal(err)
		}
		statuses = append(statuses, status)
	}
	
	
	return statuses, nil
}
func DBGetStatusStage()([]models.Stages,error){
	conn := db.GetDB()
	if conn == nil {
		return nil, errors.New("database connection is nil")
	}
	
	rows, err := conn.Query(context.Background(), "SELECT id_status_stage, name_status_stage FROM status_stages")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var statuses []models.Stages
	for rows.Next() {
		var status models.Stages
		err := rows.Scan(&status.Id_status_stage, &status.Name_status_stage)
		if err != nil {
			log.Fatal(err)
		}
		statuses = append(statuses, status)
	}
	
	return statuses, nil
}

func DBGetTypeContract()([]models.Contracts, error) {
	conn := db.GetDB()
	if conn == nil {
		return nil,errors.New("database connection is nil")
	}
	
	rows, err := conn.Query(context.Background(), "SELECT id_type_contract, name_type_contract FROM types_contracts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()	
	var types []models.Contracts
	for rows.Next() {
		var type_contract models.Contracts
		err := rows.Scan(&type_contract.Id_type, &type_contract.Name_type)
		if err != nil {
			log.Fatal(err)
		}
		types = append(types, type_contract)
	}
	
	
	return types,nil
}