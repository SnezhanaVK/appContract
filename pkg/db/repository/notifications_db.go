package db

import (
	"appContract/pkg/models"
	"database/sql"
)

func GetUsersToNotify(db *sql.DB) ([]models.Notification, error) {

	// Запрос для получения пользователей, которым нужно отправить уведомления
	query := `
		SELECT u.id_user, u.email, cn.id_notification_settings, cn.variant_notification_settings
		FROM users u
		JOIN contract_notifications cn ON u.id_user = cn.id_user
		JOIN notification_settings ns ON cn.id_notification_settings = ns.id_notification_settings
		WHERE ns.variant_notification_settings IN ('1', '3', '7')
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var  notification models.Notification
		err = rows.Scan(&notification.Id_user, &notification.Email, &notification.Id_notification_settings, &notification.Variant_notification_settings)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func GetContractInfo(db *sql.DB, contractID int) (models.Contracts, error) {
	// Запрос для получения информации о контракте
	query := `
		SELECT * FROM contracts WHERE id_contract = $1
	`
	var contract models.Contracts
	err := db.QueryRow(query, contractID).Scan(&contract.Id_contract, &contract.Name_contract, &contract.Date_contract_create, &contract.Id_user, &contract.Date_conclusion, &contract.Date_end, &contract.Id_type, &contract.Cost, &contract.Object_contract, &contract.Term_contract, &contract.Id_counterparty, &contract.Id_status_contract, &contract.Notes, &contract.Condition)
	if err != nil {
		return models.Contracts{}, err
	}
	return contract, nil
}

func GetStageInfo(db *sql.DB, stageID int) (models.Stages, error) {
	// Запрос для получения информации о этапе
	query := `
		SELECT * FROM stages WHERE id_stage = $1
	`
	var stage models.Stages
	err := db.QueryRow(query, stageID).Scan(&stage.Id_stage, &stage.Name_stage, &stage.Id_user, &stage.Description, &stage.Date_create_start, &stage.Date_create_end, &stage.Id_contract)
	if err != nil {
		return models.Stages{}, err
	}
	return stage, nil
}