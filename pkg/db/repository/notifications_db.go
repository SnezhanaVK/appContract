package db

import (
	"appContract/pkg/models"
	"database/sql"
	"time"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) GetPendingNotifications(currentDate time.Time) ([]models.PendingNotification, error) {
	// Получаем уведомления по контрактам
	contractNotifications, err := r.getContractNotifications(currentDate)
	if err != nil {
		return nil, err
	}

	// Получаем уведомления по этапам
	stageNotifications, err := r.getStageNotifications(currentDate)
	if err != nil {
		return nil, err
	}

	// Объединяем результаты
	return append(contractNotifications, stageNotifications...), nil
}

func (r *NotificationRepository) getContractNotifications(currentDate time.Time) ([]models.PendingNotification, error) {
	query := `
		SELECT 
			u.id_user, u.email, 
			CONCAT(u.surname, ' ', u.username, ' ', u.patronymic) as full_name,
			ns.id_notification_settings, ns.variant_notification_settings,
			c.id_contract, c.name_contract, c.date_end, c.notes
		FROM contract_notifications cn
		JOIN users u ON cn.id_user = u.id_user
		JOIN notification_settings ns ON cn.id_notification_settings = ns.id_notification_settings
		JOIN contracts c ON cn.id_contract = c.id_contract
		WHERE c.date_end BETWEEN $1 AND $2
	`

	// Вычисляем диапазон дат для выборки (сегодня + настройки уведомлений)
	startDate := currentDate
	endDate := currentDate.AddDate(0, 0, 7) // Максимальный период уведомления - 7 дней

	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.PendingNotification

	for rows.Next() {
		var recipient models.NotificationRecipient
		var contract models.ContractNotification
		var daysBefore int

		err := rows.Scan(
			&recipient.UserID, &recipient.Email, &recipient.FullName,
			&recipient.SettingID, &recipient.SettingVariant,
			&contract.ContractID, &contract.ContractName, &contract.EndDate, &contract.Notes,
		)
		if err != nil {
			return nil, err
		}

		// Преобразуем вариант уведомления в дни
		switch recipient.SettingVariant {
		case "1":
			daysBefore = 1
		case "3":
			daysBefore = 3
		case "7":
			daysBefore = 7
		default:
			continue // Пропускаем неизвестные варианты
		}

		// Вычисляем дату, когда нужно отправить уведомление
		notifyDate := contract.EndDate.AddDate(0, 0, -daysBefore)

		// Если сегодня день отправки уведомления
		if notifyDate.Year() == currentDate.Year() && 
		   notifyDate.Month() == currentDate.Month() && 
		   notifyDate.Day() == currentDate.Day() {
			notifications = append(notifications, models.PendingNotification{
				Recipient:  recipient,
				Contract:   &contract,
				NotifyDate: notifyDate,
			})
		}
	}

	return notifications, nil
}

func (r *NotificationRepository) getStageNotifications(currentDate time.Time) ([]models.PendingNotification, error) {
	query := `
		SELECT 
			u.id_user, u.email, 
			CONCAT(u.surname, ' ', u.username, ' ', u.patronymic) as full_name,
			ns.id_notification_settings, ns.variant_notification_settings,
			s.id_stage, s.name_stage, s.date_create_end, s.description,
			c.name_contract
		FROM stage_notifications sn
		JOIN users u ON sn.id_user = u.id_user
		JOIN notification_settings ns ON sn.id_notification_settings = ns.id_notification_settings
		JOIN stages s ON sn.id_stage = s.id_stage
		JOIN contracts c ON s.id_contract = c.id_contract
		WHERE s.date_create_end BETWEEN $1 AND $2
	`

	startDate := currentDate
	endDate := currentDate.AddDate(0, 0, 7)

	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.PendingNotification

	for rows.Next() {
		var recipient models.NotificationRecipient
		var stage models.StageNotification
		var daysBefore int

		err := rows.Scan(
			&recipient.UserID, &recipient.Email, &recipient.FullName,
			&recipient.SettingID, &recipient.SettingVariant,
			&stage.StageID, &stage.StageName, &stage.EndDate, &stage.Description,
			&stage.ContractName,
		)
		if err != nil {
			return nil, err
		}

		switch recipient.SettingVariant {
		case "1":
			daysBefore = 1
		case "3":
			daysBefore = 3
		case "7":
			daysBefore = 7
		default:
			continue
		}

		notifyDate := stage.EndDate.AddDate(0, 0, -daysBefore)

		if notifyDate.Year() == currentDate.Year() && 
		   notifyDate.Month() == currentDate.Month() && 
		   notifyDate.Day() == currentDate.Day() {
			notifications = append(notifications, models.PendingNotification{
				Recipient:  recipient,
				Stage:      &stage,
				NotifyDate: notifyDate,
			})
		}
	}

	return notifications, nil
}