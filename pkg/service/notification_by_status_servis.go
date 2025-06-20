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
	<div style="
		background-color: rgb(67, 73, 72);
		color: #ffffff;
		padding: 25px;
		font-family: 'Segoe UI', Arial, sans-serif;
		border-radius: 8px;
		max-width: 600px;
		margin: 0 auto;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
	">
		<h2 style="color: #ffffff; font-size: 18px; margin-top: 0;">Уведомление об изменении статуса</h2>
		<p style="font-size: 15px; line-height: 1.5;">Уважаемый пользователь,</p>
		<p style="font-size: 15px; line-height: 1.5;">
			Статус этапа <strong style="color: #ffffff;">%s</strong> в контракте <strong style="color: #ffffff;">%s</strong> был изменен на <strong style="color: #ffffff;">%s</strong>.
		</p>
		<div style="margin-top: 20px; padding-top: 10px; border-top: 1px solid rgba(255, 255, 255, 0.1);">
			<p style="font-size: 13px; color: #cccccc;">Это автоматическое уведомление. Пожалуйста, не отвечайте на него.</p>
		</div>
	</div>
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
