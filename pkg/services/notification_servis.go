package service

// notification_servis.go в пакете service

import (
	db "appContract/pkg/db/repository"
	"appContract/pkg/models"
	"appContract/pkg/utils"
	"fmt"
	"log"
	"time"
)

type NotificationService struct {
	repo    *db.NotificationRepository
	emailer *utils.EmailSender
}

func NewNotificationService(repo *db.NotificationRepository, emailer *utils.EmailSender) *NotificationService {
	return &NotificationService{
		repo:    repo,
		emailer: emailer,
	}
}

func (s *NotificationService) ProcessDailyNotifications() error {
    currentDate := time.Now()
    log.Printf("Начало обработки уведомлений за %s", currentDate.Format("2006-01-02"))
    
    notifications, err := s.repo.GetPendingNotifications(currentDate)
    if err != nil {
        log.Printf("Ошибка получения уведомлений: %v", err)
        return err
    }

    if len(notifications) == 0 {
        log.Println("Нет уведомлений для отправки")
        return nil
    }

    log.Printf("Найдено %d уведомлений для отправки", len(notifications))
    
    for _, n := range notifications {
        content := s.prepareEmailContent(n)
        log.Printf("Отправка уведомления для %s (%s)", n.Recipient.Email, content.Subject)
        
        if err := s.emailer.SendNotification(n.Recipient.Email, content); err != nil {
            log.Printf("Ошибка отправки уведомления %s: %v", n.Recipient.Email, err)
            continue
        }
        
        log.Printf("Уведомление отправлено: %s", n.Recipient.Email)
    }

    return nil
}

func (s *NotificationService) prepareEmailContent(n models.PendingNotification) utils.EmailContent {
	if n.Contract != nil {
		return utils.EmailContent{
			Subject: "Напоминание о контракте: " + n.Contract.ContractName,
			Body:    s.generateContractEmailBody(n.Recipient, *n.Contract),
		}
	} else {
		return utils.EmailContent{
			Subject: "Напоминание о этапе: " + n.Stage.StageName,
			Body:    s.generateStageEmailBody(n.Recipient, *n.Stage),
		}
	}
}

func (s *NotificationService) generateContractEmailBody(recipient models.NotificationRecipient, contract models.ContractNotification) string {
	daysLeft := int(time.Until(contract.EndDate).Hours() / 24)
	return `
		<html>
		<body>
			<h2>Уведомление о контракте</h2>
			<p>Уважаемый(ая) ` + recipient.FullName + `,</p>
			<p>Напоминаем вам о приближающемся сроке завершения контракта:</p>
			<ul>
				<li><strong>Контракт:</strong> ` + contract.ContractName + `</li>
				<li><strong>Дата завершения:</strong> ` + contract.EndDate.Format("02.01.2006") + `</li>
				<li><strong>Осталось дней:</strong> ` + fmt.Sprintf("%d", daysLeft) + `</li>
			</ul>
			<p>Примечания: ` + contract.Notes + `</p>
			<p>Пожалуйста, примите необходимые меры.</p>
			<br>
			<p>С уважением,<br>Система управления контрактами</p>
		</body>
		</html>
	`
}

func (s *NotificationService) generateStageEmailBody(recipient models.NotificationRecipient, stage models.StageNotification) string {
	daysLeft := int(time.Until(stage.EndDate).Hours() / 24)
	return `
		<html>
		<body>
			<h2>Уведомление о этапе</h2>
			<p>Уважаемый(ая) ` + recipient.FullName + `,</p>
			<p>Напоминаем вам о приближающемся сроке завершения этапа:</p>
			<ul>
				<li><strong>Этап:</strong> ` + stage.StageName + `</li>
				<li><strong>Контракт:</strong> ` + stage.ContractName + `</li>
				<li><strong>Дата завершения:</strong> ` + stage.EndDate.Format("02.01.2006") + `</li>
				<li><strong>Осталось дней:</strong> ` + fmt.Sprintf("%d", daysLeft) + `</li>
			</ul>
			<p>Описание: ` + stage.Description + `</p>
			<p>Пожалуйста, примите необходимые меры.</p>
			<br>
			<p>С уважением,<br>Система управления контрактами</p>
		</body>
		</html>
	`
}