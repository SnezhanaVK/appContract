package service

// import (
// 	"appContract/pkg/models"
// 	"database/sql"
// 	"errors"
// 	"net/smtp"
// 	"time"
// )

// func SendNotifications(db *sql.DB) error {
// 	users, err := db.GetUsersToNotify(db)
// 	if err != nil {
// 		return err
// 	}

// 	for _, user := range users {
// 		contract, err := db.GetContractInfo(db, user.NotificationSettingsID)
// 		if err != nil {
// 			return err
// 		}

// 		stage, err := db.GetStageInfo(db, user.NotificationSettingsID)
// 		if err != nil {
// 			return err
// 		}

// 		// Вычислить количество дней для вычитания из текущей даты
// 		var daysToSubtract int
// 		switch user.VariantNotificationSettings {
// 		case "1":
// 			daysToSubtract = 1
// 		case "3":
// 			daysToSubtract = 3
// 		case "7":
// 			daysToSubtract = 7
// 		default:
// 			return errors.New("неверный вариант уведомления")
// 		}

// 		// Проверить, если сегодня день отправки уведомления
// 		if time.Now().Day() == contract.DateEnd.Day()-daysToSubtract {
// 			// Отправить уведомление
// 			err = sendNotification(user.Email, contract, stage)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// func sendNotification(email string, contract models.Contract, stage models.Stage) error {
// 	// Используйте SMTP для отправки уведомлений
// 	auth := smtp.PlainAuth("", "your_email", "your_password", "smtp.gmail.com")

// 	// Создайте сообщение
// 	msg := "Уведомление о завершении этапа/контракта"
// 	to := []string{email}

// 	// Отправьте уведомление
// 	err := smtp.SendMail("smtp.gmail.com:587", auth, "your_email", to, []byte(msg))
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func sendNotification(email string, variant_notification_settings string, contract models.Contracts, stage models.Stages) error {
//     // Используйте SMTP для отправки уведомлений
//     auth := smtp.PlainAuth("", "your_email", "your_password", "smtp.gmail.com")

//     // Создайте сообщение
//     msg := "Уведомление о завершении этапа/контракта"
// 	to := []string{email}

//     // Отправьте уведомление
//     err := smtp.SendMail("smtp.gmail.com:587", auth, "your_email", to, []byte(msg))
//     if err != nil {
//         return err
//     }

//     return nil
// }