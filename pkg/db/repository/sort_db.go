package db

//sort_db.go
import (
	db "appContract/pkg/db"
	"appContract/pkg/models"
	"context"
	"errors"
	"fmt"
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
			&tag.Tags_contract)
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
		var status models.Contracts
		err := rows.Scan(&status.Id_status_contract, &status.Name_status_contract)
		if err != nil {
			log.Fatal(err)
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}
func DBGetStatusStage() ([]models.Stages, error) {
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

func DBGetTypeContract() ([]models.Contracts, error) {
	conn := db.GetDB()
	if conn == nil {
		return nil, errors.New("database connection is nil")
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

	return types, nil
}

func AddTagToContract(contractID int, tagID int) error {
	conn := db.GetDB()
	if conn == nil {
		return errors.New("connection error")
	}
	var exists bool
	err := conn.QueryRow(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM tegs WHERE id_teg = $1)`,
		tagID).Scan(&exists)

	if err != nil {
		return fmt.Errorf("tag check error: %v", err)
	}

	if !exists {
		return errors.New("tag does not exist")
	}
	_, err = conn.Exec(context.Background(),
		`INSERT INTO contracts_by_tegs (id_contract, id_teg)
         VALUES ($1, $2)
         ON CONFLICT (id_contract, id_teg) DO NOTHING`,
		contractID, tagID)

	if err != nil {
		return fmt.Errorf("failed to add tag: %v", err)
	}

	return nil
}
func RemoveTagFromContract(contractID int, tagID int) error {
	conn := db.GetDB()
	if conn == nil {
		return errors.New("connection error")
	}

	result, err := conn.Exec(context.Background(),
		`DELETE FROM contracts_by_tegs 
         WHERE id_contract = $1 AND id_teg = $2`,
		contractID, tagID)

	if err != nil {
		return fmt.Errorf("failed to remove tag: %v", err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return errors.New("tag association not found")
	}

	return nil
}

type Tag struct {
	ID   int
	Name string
}

func GetContractIDTags(contractID int) ([]Tag, error) {
	conn := db.GetDB()
	if conn == nil {
		return nil, errors.New("connection error")
	}
	rows, err := conn.Query(context.Background(),
		`SELECT t.id_teg, t.name_teg 
         FROM contracts_by_tegs cbt
         JOIN tegs t ON cbt.id_teg = t.id_teg
         WHERE cbt.id_contract = $1`,
		contractID)

	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var tags []Tag

	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return tags, nil
}
