// db/notifications_db.go
package db

import (
	db "appContract/pkg/db"
	"appContract/pkg/models"
	"context"
	"errors"
	"fmt"
)

func GetContractNotifications() ([]models.ContractNotification, error) {
	conn := db.GetDB()

	query := `
		SELECT 
			u.email, 
			c.date_end, 
			ns.variant_notification_settings
		FROM users u
		JOIN contracts c ON u.id_user = c.id_user
		JOIN notification_settings_by_user nsu ON u.id_user = nsu.id_user
		JOIN notification_settings ns ON nsu.id_notification_settings = ns.id_notification_settings
		WHERE c.date_end - ns.variant_notification_settings = CURRENT_DATE`

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("error querying contract notifications: %v", err)
	}
	defer rows.Close()

	var notifications []models.ContractNotification
	for rows.Next() {
		var n models.ContractNotification
		err := rows.Scan(&n.Email, &n.DateEnd, &n.DaysBefore)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		notifications = append(notifications, n)
	}

	return notifications, rows.Err()
}

func GetStageNotifications() ([]models.StageNotification, error) {
	conn := db.GetDB()

	query := `
		SELECT 
			u.email, 
			s.date_create_end, 
			ns.variant_notification_settings
		FROM users u
		JOIN stages s ON u.id_user = s.id_user
		JOIN notification_settings_by_user nsu ON u.id_user = nsu.id_user
		JOIN notification_settings ns ON nsu.id_notification_settings = ns.id_notification_settings
		WHERE s.date_create_end - ns.variant_notification_settings = CURRENT_DATE`

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("error querying stage notifications: %v", err)
	}
	defer rows.Close()

	var notifications []models.StageNotification
	for rows.Next() {
		var n models.StageNotification
		err := rows.Scan(&n.Email, &n.DateEnd, &n.DaysBefore)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		notifications = append(notifications, n)
	}

	return notifications, rows.Err()
}

func SetUserNotificationSettings(userID int, variants []int) error {
	conn := db.GetDB()
	if conn == nil {
		return errors.New("connection error")
	}

	// Начинаем транзакцию
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("transaction start error: %v", err)
	}
	defer tx.Rollback(context.Background())

	// Удаляем старые настройки
	_, err = tx.Exec(context.Background(),
		`DELETE FROM notification_settings_by_user 
		WHERE id_user = $1`,
		userID)
	if err != nil {
		return fmt.Errorf("failed to delete old settings: %v", err)
	}

	// Если передан пустой список - коммитим и выходим
	if len(variants) == 0 {
		return tx.Commit(context.Background())
	}

	// Добавляем новые настройки
	query := `
	INSERT INTO notification_settings_by_user 
		(id_user, id_notification_settings)
	SELECT $1, ns.id_notification_settings 
	FROM notification_settings ns
	WHERE ns.variant_notification_settings = ANY($2)`

	_, err = tx.Exec(context.Background(), query, userID, variants)
	if err != nil {
		return fmt.Errorf("failed to insert new settings: %v", err)
	}

	return tx.Commit(context.Background())
}

func GetUserNotificationSettings(userID int) ([]int, error) {
	conn := db.GetDB()
	if conn == nil {
		return nil, errors.New("connection error")
	}

	query := `
	SELECT ns.variant_notification_settings 
	FROM notification_settings_by_user nsu
	JOIN notification_settings ns 
		ON ns.id_notification_settings = nsu.id_notification_settings
	WHERE nsu.id_user = $1`

	rows, err := conn.Query(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %v", err)
	}
	defer rows.Close()

	var variants []int
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		variants = append(variants, v)
	}

	return variants, nil
}
