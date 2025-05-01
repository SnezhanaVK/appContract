// db/notifications_db.go
package db

// notifications_db.go в пакете db
import (
	db "appContract/pkg/db"
	"appContract/pkg/models"
	"context"
	"fmt"
)

func GetContractNotifications() ([]models.ContractNotification, error) {
	conn := db.GetDB() // Используем глобальное подключение
    
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

	rows, err := conn.Query( context.Background(),query)
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
	conn :=db.GetDB() // Используем глобальное подключение
    
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

	rows, err := conn.Query( context.Background(),query)
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