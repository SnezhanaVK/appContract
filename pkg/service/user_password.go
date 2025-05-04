package service

import (
	"fmt"
	"log"

	db "appContract/pkg/db/repository"
	"appContract/pkg/models"
	"appContract/pkg/utils"
)

// CreateUser создает пользователя с автоматически сгенерированным паролем
func CreateUser(user models.Users, emailSender utils.EmailSenderInterface) error {
	// Генерируем постоянный пароль
	password, err := utils.GenerateStrongPassword()
	if err != nil {
		return fmt.Errorf("failed to generate password: %v", err)
	}

	// Добавляем пользователя в БД
	if err := db.DBaddUser(user, password); err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	// Отправляем письмо с постоянным паролем
	emailContent := utils.EmailContent{
		Subject: "Данные вашего аккаунта",
		Body: fmt.Sprintf(`
			<h2>Учетные данные для входа</h2>
			<p>Здравствуйте, %s %s %s!</p>
			<p>Ваш аккаунт был успешно создан. Сохраните эти данные:</p>
			<p><strong>Логин:</strong> %s</p>
			<p><strong>Пароль:</strong> %s</p>
			<p style="color: #ff0000; font-weight: bold;">
				Это ваш постоянный пароль. Никому его не сообщайте!
			</p>
			<p>Для входа перейдите по <a href="https://ваш-сайт/login">ссылке</a>.</p>
		`, user.Surname, user.Username, user.Patronymic, user.Login, password),
	}

	if err := emailSender.SendNotification(user.Email, emailContent); err != nil {
		log.Printf("Failed to send email to %s: %v", user.Email, err)
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
