package service

import (
	"fmt"
	"log"

	db "appContract/pkg/db/repository"
	"appContract/pkg/models"
	"appContract/pkg/utils"
)

func CreateUser(user models.Users, emailSender utils.EmailSenderInterface)  error {

	password, err := utils.GenerateStrongPassword()
	if err != nil {
		return fmt.Errorf("failed to generate password: %v", err)
	}

	if err := db.DBaddUser(user, password); err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	emailContent := utils.EmailContent{
    Subject: "Данные вашего аккаунта",
    Body: fmt.Sprintf(`
        <div style="
            background-color: #ffffff;
            color: #333333;
            padding: 25px;
            font-family: 'Segoe UI', Arial, sans-serif;
            border-radius: 8px;
            max-width: 600px;
            margin: 0 auto;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            border: 1px solid #e0e0e0;
        ">
            <h2 style="color: #333333; font-size: 18px; margin-top: 0;">Учетные данные для входа</h2>
            <p style="font-size: 15px; line-height: 1.5;">Здравствуйте, <strong style="color: #333333;">%s %s %s</strong>!</p>
            <p style="font-size: 15px; line-height: 1.5;">Ваш аккаунт был успешно создан. Сохраните эти данные:</p>
            
            <div style="
                background: #f8f9fa;
                padding: 15px;
                border-radius: 6px;
                margin: 15px 0;
                border: 1px solid #e0e0e0;
            ">
                <p style="font-size: 15px; line-height: 1.5;">
                    <strong style="display: inline-block; width: 80px; color: #333333;">Логин:</strong> 
                    <span style="font-family: monospace; color: #333333;">%s</span>
                </p>
                <p style="font-size: 15px; line-height: 1.5;">
                    <strong style="display: inline-block; width: 80px; color: #333333;">Пароль:</strong> 
                    <span style="
                        font-family: monospace;
                        color: #d32f2f;
                        font-weight: bold;
                    ">%s</span>
                </p>
            </div>
            
            <p style="
                font-size: 14px;
                line-height: 1.5;
                color: #d32f2f;
                font-weight: bold;
                padding: 10px;
                background: #ffebee;
                border-left: 3px solid #d32f2f;
                margin: 15px 0;
            ">
                Это ваш постоянный пароль. Никому его не сообщайте!
            </p>
            
            <p style="font-size: 15px; line-height: 1.5; color: #333333;">
                Для входа перейдите по 
                <a href="https://ваш-сайт/login" style="
                    color: #1976d2;
                    text-decoration: none;
                    font-weight: bold;
                ">ссылке</a>.
            </p>
            
            <div style="margin-top: 20px; padding-top: 10px; border-top: 1px solid #e0e0e0;">
                <p style="font-size: 13px; color: #757575;">Это автоматическое уведомление. Пожалуйста, не отвечайте на него.</p>
            </div>
        </div>
    `, user.Surname, user.Username, user.Patronymic, user.Login, password),
}

	if err := emailSender.SendNotification(user.Email, emailContent); err != nil {
		log.Printf("Failed to send email to %s: %v", user.Email, err)
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
