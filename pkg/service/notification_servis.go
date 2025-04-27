// service/notification_service.go
package service

import (
	db "appContract/pkg/db/repository"
	"appContract/pkg/models"
	"appContract/pkg/utils"
	"fmt"
	"log"
)
var emailSender *utils.EmailSender // Глобальная переменная
func InitEmailSender(es *utils.EmailSender) {
	emailSender = es
}
func ProcessDailyNotifications() error {
    contracts, err := db.GetContractNotifications()
    if err != nil {
        return err
    }

    stages, err := db.GetStageNotifications()
    if err != nil {
        return err
    }

    for _, cn := range contracts {
        content := utils.EmailContent{
            Subject: "Contract Expiration Notification",
            Body:    formatContractMessage(cn),
        }
        if err := emailSender.SendNotification(cn.Email, content); err != nil {
            log.Printf("Failed to send notification to %s: %v", cn.Email, err)
        }
    }

    for _, sn := range stages {
        content := utils.EmailContent{
            Subject: "Stage Deadline Notification",
            Body:    formatStageMessage(sn),
        }
        if err := emailSender.SendNotification(sn.Email, content); err != nil {
            log.Printf("Failed to send notification to %s: %v", sn.Email, err)
        }
    }

    return nil
}

func formatContractMessage(cn models.ContractNotification) string {
    return fmt.Sprintf(`
        <p>Your contract ending on %s will expire in %d days.</p>
        <p>Please take necessary actions.</p>
    `, cn.DateEnd.Format("2006-01-02"), cn.DaysBefore)
}

func formatStageMessage(sn models.StageNotification) string {
    return fmt.Sprintf(`
        <p>Your stage ending on %s will expire in %d days.</p>
        <p>Please check the project timeline.</p>
    `, sn.DateEnd.Format("2006-01-02"), sn.DaysBefore)
}