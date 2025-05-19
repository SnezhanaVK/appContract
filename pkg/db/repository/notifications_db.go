// db/notifications_db.go
package db

import (
	db "appContract/pkg/db"
	"appContract/pkg/models"
	"context"
	"errors"
	"fmt"
	"log"
)

func GetContractNotifications() ([]models.ContractNotification, error) {
	conn := db.GetDB()

	query := `
		SELECT 
			u.id_user,
			u.email, 
			c.id_contract,
			c.name_contract,
			c.date_end, 
			ns.variant_notification_settings
		FROM users u
		JOIN contracts c ON u.id_user = c.id_user
		JOIN notification_settings_by_user nsu ON u.id_user = nsu.id_user
		JOIN notification_settings ns ON nsu.id_notification_settings = ns.id_notification_settings
		WHERE (c.date_end - CURRENT_DATE) = ns.variant_notification_settings
		AND c.date_end >= CURRENT_DATE`

	log.Printf("[БД] Ищем контракты для уведомлений")
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса контрактов: %w", err)
	}
	defer rows.Close()

	var notifications []models.ContractNotification
	for rows.Next() {
		var n models.ContractNotification
		err := rows.Scan(
			&n.UserID,
			&n.Email,
			&n.ContractID,
			&n.ContractName,
			&n.DateEnd,
			&n.DaysBefore,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования: %w", err)
		}
		notifications = append(notifications, n)
	}

	log.Printf("[БД] Найдено %d контрактов для уведомления", len(notifications))
	return notifications, rows.Err()
}

func GetStageNotifications() ([]models.StageNotification, error) {
	conn := db.GetDB()

	query := `
		SELECT 
			u.id_user,
			u.email,
			s.id_stage,
			s.name_stage,
			s.id_contract,
			s.date_create_end, 
			ns.variant_notification_settings
		FROM users u
		JOIN stages s ON u.id_user = s.id_user
		JOIN notification_settings_by_user nsu ON u.id_user = nsu.id_user
		JOIN notification_settings ns ON nsu.id_notification_settings = ns.id_notification_settings
		WHERE (s.date_create_end - CURRENT_DATE) = ns.variant_notification_settings
		AND s.date_create_end >= CURRENT_DATE`

	log.Printf("[БД] Ищем этапы для уведомлений")
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса этапов: %w", err)
	}
	defer rows.Close()

	var notifications []models.StageNotification
	for rows.Next() {
		var n models.StageNotification
		err := rows.Scan(
			&n.UserID,
			&n.Email,
			&n.StageID,
			&n.StageName,
			&n.ContractID,
			&n.DateEnd,
			&n.DaysBefore,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования: %w", err)
		}
		notifications = append(notifications, n)
	}

	log.Printf("[БД] Найдено %d этапов для уведомления", len(notifications))
	return notifications, rows.Err()
}
func SetUserNotificationSettings(userID int, variants []int) error {
	conn := db.GetDB()
	if conn == nil {
		return errors.New("connection error")
	}
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("transaction start error: %v", err)
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(context.Background(),
		`DELETE FROM notification_settings_by_user 
		WHERE id_user = $1`,
		userID)
	if err != nil {
		return fmt.Errorf("failed to delete old settings: %v", err)
	}
	if len(variants) == 0 {
		return tx.Commit(context.Background())
	}
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
