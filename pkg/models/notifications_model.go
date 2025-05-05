package models

import "time"

type NotificationSettings struct {
	ID   int `json:"id_notification_settings"`
	Days int `json:"variant_notification_settings"`
}

type UserNotification struct {
	UserID     int    `json:"id_user"`
	Email      string `json:"email"`
	DaysBefore int    `json:"days_before"`
}

type ContractNotification struct {
	UserNotification
	ContractID   int       `json:"id_contract"`
	ContractName string    `json:"name_contract"`
	DateEnd      time.Time `json:"date_end"`
}

type StageNotification struct {
	UserNotification
	StageID    int       `json:"id_stage"`
	StageName  string    `json:"name_stage"`
	ContractID int       `json:"id_contract"`
	DateEnd    time.Time `json:"date_create_end"`
}
