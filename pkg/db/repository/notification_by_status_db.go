package db

import (
	"appContract/pkg/db"
	"context"
	"fmt"
)

// GetUsersToNotifyForStage возвращает пользователей, которых нужно уведомить об изменении этапа
func GetUsersToNotifyForStage(stageID int) ([]int, error) {
	coon := db.GetDB
	var userIDs []int

	query := `
		SELECT DISTINCT u.id_user
		FROM users u
		JOIN user_by_role ur ON u.id_user = ur.id_user
		JOIN contracts c ON c.id_user = u.id_user
		JOIN stages s ON s.id_contract = c.id_contract
		WHERE s.id_stage = $1
		UNION
		SELECT s.id_user
		FROM stages s
		WHERE s.id_stage = $1
	`

	rows, err := coon().Query(context.Background(), query, stageID)
	if err != nil {
		return nil, fmt.Errorf("error querying users to notify: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("error scanning user id: %v", err)
		}
		userIDs = append(userIDs, id)
	}

	return userIDs, nil
}

// GetUserEmail возвращает email пользователя по ID
func GetUserEmail(userID int) (string, error) {
	coon := db.GetDB
	var email string
	err := coon().QueryRow(context.Background(), "SELECT email FROM users WHERE id_user = $1", userID).Scan(&email)
	if err != nil {
		return "", fmt.Errorf("error getting user email: %v", err)
	}
	return email, nil
}

// GetStageInfo возвращает информацию об этапе для уведомления
func GetStageInfo(stageID int) (stageName, contractName string, err error) {
	coon := db.GetDB
	query := `
		SELECT s.name_stage, c.name_contract
		FROM stages s
		JOIN contracts c ON s.id_contract = c.id_contract
		WHERE s.id_stage = $1
	`
	err = coon().QueryRow(context.Background(), query, stageID).Scan(&stageName, &contractName)
	if err != nil {
		return "", "", fmt.Errorf("error getting stage info: %v", err)
	}
	return stageName, contractName, nil
}

// GetStatusName возвращает название статуса по ID
func GetStatusName(statusID int) (string, error) {
	coon := db.GetDB
	var name string
	err := coon().QueryRow(context.Background(), "SELECT name_status_stage FROM status_stages WHERE id_status_stage = $1", statusID).Scan(&name)
	if err != nil {
		return "", fmt.Errorf("error getting status name: %v", err)
	}
	return name, nil
}
