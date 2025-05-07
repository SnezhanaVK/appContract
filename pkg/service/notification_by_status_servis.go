package service

import (
	db "appContract/pkg/db/repository"
	"appContract/pkg/utils"
	"context"
	"fmt"
	"log"
)

type NotificationType string

const (
	StageStatusChange NotificationType = "stage_status_change"
)

type NotificationService struct {
	mailer utils.EmailSenderInterface
}

func NewNotificationService(mailer utils.EmailSenderInterface) *NotificationService {
	return &NotificationService{mailer: mailer}
}

func (s *NotificationService) NotifyStageStatusChange(ctx context.Context, stageID, statusID int, comment string) error {
	log.Printf("Starting notification process for stage %d, status %d", stageID, statusID)

	userIDs, err := db.GetUsersToNotifyForStage(stageID)
	if err != nil {
		log.Printf("Error getting users to notify: %v", err)
		return fmt.Errorf("error getting users to notify: %v", err)
	}
	log.Printf("Found %d users to notify", len(userIDs))

	stageName, contractName, err := db.GetStageInfo(stageID)
	if err != nil {
		log.Printf("Error getting stage info: %v", err)
		return fmt.Errorf("error getting stage info: %v", err)
	}

	statusName, err := db.GetStatusName(statusID)
	if err != nil {
		log.Printf("Error getting status name: %v", err)
		return fmt.Errorf("error getting status name: %v", err)
	}

	subject := fmt.Sprintf("Изменение статуса этапа '%s' в контракте '%s'", stageName, contractName)
	body := fmt.Sprintf(`
		<h1>Уведомление об изменении статуса</h1>
		<p>Статус этапа <strong>%s</strong> в контракте <strong>%s</strong> был изменен на <strong>%s</strong>.</p>
	`, stageName, contractName, statusName)

	if comment != "" {
		body += fmt.Sprintf(`<p>Комментарий: <em>%s</em></p>`, comment)
	}

	for _, userID := range userIDs {
		email, err := db.GetUserEmail(userID)
		if err != nil {
			log.Printf("Error getting email for user %d: %v", userID, err)
			continue
		}

		log.Printf("Preparing to send email to %s", email)

		content := utils.EmailContent{
			Subject: subject,
			Body:    body,
		}

		if err := s.mailer.SendNotification(email, content); err != nil {
			log.Printf("Error sending email to %s: %v", email, err)
		} else {
			log.Printf("Successfully sent notification to %s", email)
		}
	}

	return nil
}
