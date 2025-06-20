package service

import (
	db "appContract/pkg/db/repository"
	"time"

	"appContract/pkg/models"
	"appContract/pkg/utils"
	"fmt"
	"log"
)

var emailSender *utils.EmailSender

func InitEmailSender(es *utils.EmailSender) {
	emailSender = es
}

func ProcessDailyNotifications() error {
	log.Println("=== Начало обработки уведомлений ===")

	currentDate := time.Now().Format("2006-01-02")
	log.Printf("Текущая дата: %s", currentDate)

	contracts, err := db.GetContractNotifications()
	if err != nil {
		log.Printf("Ошибка получения контрактов: %v", err)
		return err
	}
	log.Printf("Найдено контрактов для уведомления: %d", len(contracts))

	stages, err := db.GetStageNotifications()
	if err != nil {
		log.Printf("Ошибка получения этапов: %v", err)
		return err
	}
	log.Printf("Найдено этапов для уведомления: %d", len(stages))

	sendNotifications(contracts, stages)

	log.Println("=== Обработка завершена ===")
	return nil
}
func formatContractMessage(cn models.ContractNotification) string {
	daysText := formatDaysText(cn.DaysBefore)
	return fmt.Sprintf(`
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
			<h2 style="color: #ffffff; font-size: 18px; margin-top: 0;">Уведомление о контракте</h2>
			<p style="font-size: 15px; line-height: 1.5;">Уважаемый пользователь,</p>
			<p style="font-size: 15px; line-height: 1.5;">
				Контракт <strong style="color: #ffffff;">%s</strong> завершается через <strong>%s</strong>.
			</p>
			<p style="font-size: 15px; line-height: 1.5;">Дата окончания: <strong>%s</strong></p>
			<p style="font-size: 15px; line-height: 1.5;">Пожалуйста, примите необходимые меры.</p>
			<div style="margin-top: 20px; padding-top: 10px; border-top: 1px solid rgba(255, 255, 255, 0.1);">
				<p style="font-size: 13px; color: #cccccc;">Это автоматическое уведомление. Пожалуйста, не отвечайте на него.</p>
			</div>
		</div>
	`, cn.ContractName, daysText, cn.DateEnd.Format("02.01.2006"))
}

func formatStageMessage(sn models.StageNotification) string {
	daysText := formatDaysText(sn.DaysBefore)
	return fmt.Sprintf(`
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
			<h2 style="color: #ffffff; font-size: 18px; margin-top: 0;">Уведомление об этапе</h2>
			<p style="font-size: 15px; line-height: 1.5;">Уважаемый пользователь,</p>
			<p style="font-size: 15px; line-height: 1.5;">
				Этап <strong style="color: #ffffff;">%s</strong> завершается через <strong>%s</strong>.
			</p>
			<p style="font-size: 15px; line-height: 1.5;">Дата окончания: <strong>%s</strong></p>
			<p style="font-size: 15px; line-height: 1.5;">Пожалуйста, проверьте сроки выполнения.</p>
			<div style="margin-top: 20px; padding-top: 10px; border-top: 1px solid rgba(255, 255, 255, 0.1);">
				<p style="font-size: 13px; color: #cccccc;">Это автоматическое уведомление. Пожалуйста, не отвечайте на него.</p>
			</div>
		</div>
	`, sn.StageName, daysText, sn.DateEnd.Format("02.01.2006"))
}


func formatDaysText(days int) string {
	switch {
	case days == 1:
		return "1 день"
	case days >= 2 && days <= 4:
		return fmt.Sprintf("%d дня", days)
	default:
		return fmt.Sprintf("%d дней", days)
	}
}

func sendNotifications(contracts []models.ContractNotification, stages []models.StageNotification) {
	for _, cn := range contracts {
		content := utils.EmailContent{
			Subject: fmt.Sprintf("Завершение контракта '%s'", cn.ContractName),
			Body:    formatContractMessage(cn),
		}

		log.Printf("Отправка уведомления для контракта %d (%s) пользователю %s",
			cn.ContractID, cn.ContractName, cn.Email)

		if err := emailSender.SendNotification(cn.Email, content); err != nil {
			log.Printf("Ошибка отправки: %v", err)
		}
	}

	for _, sn := range stages {
		content := utils.EmailContent{
			Subject: fmt.Sprintf("Завершение этапа '%s'", sn.StageName),
			Body:    formatStageMessage(sn),
		}

		log.Printf("Отправка уведомления для этапа %d (%s) пользователю %s",
			sn.StageID, sn.StageName, sn.Email)

		if err := emailSender.SendNotification(sn.Email, content); err != nil {
			log.Printf("Ошибка отправки: %v", err)
		}
	}
}
