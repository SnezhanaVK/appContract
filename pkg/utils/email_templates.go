package utils

const UserCreatedTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f9fa; padding: 20px; text-align: center; }
        .content { padding: 20px; }
        .credentials { background-color: #f1f1f1; padding: 15px; border-radius: 5px; }
        .footer { margin-top: 20px; font-size: 12px; color: #6c757d; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>Ваш аккаунт был создан</h2>
        </div>
        <div class="content">
            <p>Здравствуйте, %s!</p>
            <p>Ваш аккаунт в системе был успешно создан.</p>
            
            <div class="credentials">
                <p><strong>Логин:</strong> %s</p>
                <p><strong>Пароль:</strong> %s</p>
            </div>
            
            <p>Пожалуйста, измените пароль после первого входа в систему.</p>
        </div>
        <div class="footer">
            <p>Это письмо отправлено автоматически, пожалуйста, не отвечайте на него.</p>
        </div>
    </div>
</body>
</html>
`
