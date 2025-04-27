package models

import "time"

//notification_model.go в папке models



type Notification_settings struct {
    Id_notification_settings      int       `json:"id_notification"`
    Variant_notification_settings int       `json:"variant_notification"` // Изменено на int
    Id_user                       int       `json:"id_user"`
    Email                         string    `json:"email"`
    Date_end                      time.Time `json:"date_end"`          // Для контракта
    Date_create_end               time.Time `json:"date_create_end"`  // Для этапа
}

// Дополнительные структуры для запросов
type ContractNotification struct {
    Email      string    `json:"email"`
    DateEnd    time.Time `json:"date_end"`
    DaysBefore int       `json:"days_before"`
}

type StageNotification struct {
    Email      string    `json:"email"`
    DateEnd    time.Time `json:"date_create_end"`
    DaysBefore int       `json:"days_before"`
}