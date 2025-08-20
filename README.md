Contract Management System Backend
Backend-сервис для управления контрактами с системой уведомлений, ролевой моделью и workflow этапов.

 Технологии
Язык: Go (Golang) 1.24.1

Фреймворк: Стандартная библиотека HTTP + Gorilla WebSocket

База данных: PostgreSQL с пулом соединений pgx

Аутентификация: JWT токены с cookie-хранилищем

Расписания: Cron для автоматических задач

Уведомления: Email рассылка через SMTP

WebSocket: Real-time поиск

Мониторинг: Prometheus metrics

Контейнеризация: Docker


Структура проекта
appContract/
├── .vscode/                 # Конфигурация VS Code
├── cmd/                     # Точка входа (main.go)
├── pkg/                     # Основные пакеты
│   ├── db/                  # Работа с базой данных
│   │   └── repository/      # Репозитории и запросы
│   ├── handlers/            # HTTP обработчики
│   ├── middleware/          # Промежуточное ПО (CORS, Monitoring)
│   ├── models/              # Структуры данных
│   ├── routers/             # Маршрутизация
│   ├── service/             # Бизнес-логика
│   ├── unit_tests/          # Юнит-тесты
│   └── utils/               # Вспомогательные утилиты
└── Dockerfile              # Конфигурация Docker

Модель базы данных
Система включает 18 связанных таблиц:
*Пользователи с ролями (админ, менеджер) и фотографиями
*Контракты с типами, статусами и тегами
*Контрагенты с полными реквизитами (ИНН, ОГРН, адрес)
*Этапы контрактов с историей статусов
*Комментарии и файлы к этапам
*Система уведомлений с настройками пользователей
*Теги для категоризации контрактов

Функциональность

Аутентификация и авторизация
  JWT-аутентификация с http-only cookies
  Ролевая модель (Admin/Manager) с гибким назначением
  Восстановление пароля через email с верификацией кодов
  Защищенные пароли с bcrypt + salt

Управление контрактами
  Полный CRUD для контрактов
  Фильтрация по типу, статусу, тегам, дате
  Поиск через WebSocket в реальном времени
  Система этапов с трекингом статусов
  Комментарии и прикрепление файлов к этапам

Система уведомлений
 
  Ежедневные уведомления о завершении контрактов/этапов
  Email рассылка с HTML-шаблонами
  Настраиваемые уведомления для пользователей
  Верификация кодов с временным хранилищем

Дополнительные возможности

  Prometheus metrics для мониторинга
  CORS для кросс-доменных запросов
  Connection pooling для PostgreSQL
  Docker контейнеризация


  Настройка окружения
  # Обязательные переменные окружения
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=1234
DB_NAME=contract_db
SSL_MODE=disable

# Email настройки (для уведомлений)
EMAIL_FROM=your-email@example.com
EMAIL_PASSWORD=your-app-password
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587

  Запуск проекта
Локстановка зависимостей
go mod tidy
# Запуск сервера
go run cmd/main.go
# Сервер доступен по
# http://localhost:8080
# Metrics: http://localhost:8080/metrics

Docker сборка
# Сборка образа
docker build -t appcontract .
# Запуск контейнера
docker run -p 8080:8080 -e DB_HOST=host.docker.internal appcontract

 API Endpoints
Аутентификация
POST /api/authorizations - Вход в систему

GET /api/authorizations/token - Верификация токена

GET /api/authorizations/logout - Выход

PUT /api/authorizations/forgot-password - Смена пароля

POST /api/authorizations/sendingCode - Отправка кода подтверждения

POST /api/authorizations/verifyCode - Проверка кода

Пользователи
GET /api/users - Получить всех пользователей

GET /api/users/{userID} - Получить пользователя по ID

POST /api/users/create - Создать пользователя

GET /api/users/rolesUser/{userID} - Получить роли пользователя

POST /api/users/addRoleAdmin/{userID} - Добавить роль админа

POST /api/users/addRoleManager/{userID} - Добавить роль менеджера

DELETE /api/users/deleteRoleUser/{userID} - Удалить роль пользователя

DELETE /api/users/deleteRoleManager/{userID} - Удалить роль менеджера

PUT /api/users/update/{userID} - Обновить пользователя

DELETE /api/users/{id} - Удалить пользователя

Контракты
GET /api/contracts - Все контракты

GET /api/contracts/user/{userID} - Контракты пользователя

GET /api/contracts/{contractID} - Контракт по ID

GET /api/contracts/byType/{idType} - Контракты по типу

POST /api/contracts/byDateCreate - Контракты по дате

GET /api/contracts/byTeg/{id_teg_contract} - Контракты по тегу

GET /api/contracts/byStatus/{id_status_contract} - Контракты по статусу

POST /api/contracts/create - Создать контракт

PUT /api/contracts/{contractID} - Изменить контракт

PUT /api/contracts/userchange - Изменить пользователя контракта

DELETE /api/contracts/{contractID} - Удалить контракт

Уведомления
GET /api/users/{userID}/notifications - Получить настройки

PUT /api/users/{userID}/notifications - Обновить настройки

DELETE /api/users/{userID}/notifications - Удалить настройки

Дополнительные endpoints
POST /api/search - Поиск через WebSocket

GET /test-notifications - Тест уведомлений

GET /metrics - Prometheus metrics

 Безопасность
HTTPS-only cookies для JWT токенов
BCrypt с солью для хеширования паролей
Валидация входных данных на всех уровнях
CORS политики для контроля доступа
Connection pooling с таймаутами

  Мониторинг
Встроенная поддержка Prometheus metrics:
HTTP запросы/ответы
Время выполнения
Статус коды
Database connection stats

   Утилиты
Генерация надежных паролей с требованиями безопасности
Email отправка с TLS и таймаутами
Верификация кодов с автоматической очисткой
Миграции базы данных при запуске

