package db

// notifications_db.go в пакете db
import (
	"appContract/pkg/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository struct {
    db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
    return &NotificationRepository{db: db}
}

func (r *NotificationRepository) GetPendingNotifications(currentDate time.Time) ([]models.PendingNotification, error) {
    contractNotifications, err := r.getContractNotifications(currentDate)
    if err != nil {
        return nil, err
    }

    stageNotifications, err := r.getStageNotifications(currentDate)
    if err != nil {
        return nil, err
    }

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
        WHERE c.date_end >= $1
    `

    rows, err := r.db.Query(context.Background(), query, currentDate)
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

        notifyDate := contract.EndDate.AddDate(0, 0, -daysBefore)
        if isSameDay(notifyDate, currentDate) {
            notifications = append(notifications, models.PendingNotification{
                Recipient:  recipient,
                Contract:   &contract,
                NotifyDate: notifyDate,
            })
        }
    }

    return notifications, rows.Err()
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

    rows, err := r.db.Query(context.Background(), query, startDate, endDate)
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
        if isSameDay(notifyDate, currentDate) {
            notifications = append(notifications, models.PendingNotification{
                Recipient:  recipient,
                Stage:      &stage,
                NotifyDate: notifyDate,
            })
        }
    }

    return notifications, rows.Err()
}

func isSameDay(date1, date2 time.Time) bool {
    y1, m1, d1 := date1.Date()
    y2, m2, d2 := date2.Date()
    return y1 == y2 && m1 == m2 && d1 == d2
}