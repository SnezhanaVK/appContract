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
		<p>Уважаемый пользователь,</p>
		<p>Контракт <strong>%s</strong> завершается через %s.</p>
		<p>Дата окончания: %s</p>
		<p>Пожалуйста, примите необходимые меры.</p>
	`, cn.ContractName, daysText, cn.DateEnd.Format("02.01.2006"))
}

func formatStageMessage(sn models.StageNotification) string {
	daysText := formatDaysText(sn.DaysBefore)
	return fmt.Sprintf(`
		<p>Уважаемый пользователь,</p>
		<p>Этап <strong>%s</strong> завершается через %s.</p>
		<p>Дата окончания: %s</p>
		<p>Пожалуйста, проверьте сроки выполнения.</p>
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
