package models

//notification_model.go в папке models

import (
	"time"
)

type NotificationSetting struct {
	ID          int    `json:"id"`
	Variant     string `json:"variant"` // "1", "3", "7" days before
	Description string `json:"description"`
}

type NotificationRecipient struct {
	UserID      int    `json:"user_id"`
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	SettingID   int    `json:"setting_id"`
	SettingVariant string `json:"setting_variant"`
}

type ContractNotification struct {
	ContractID   int       `json:"contract_id"`
	ContractName string    `json:"contract_name"`
	EndDate      time.Time `json:"end_date"`
	Notes        string    `json:"notes"`
}

type StageNotification struct {
	StageID      int       `json:"stage_id"`
	StageName    string    `json:"stage_name"`
	EndDate      time.Time `json:"end_date"`
	Description  string    `json:"description"`
	ContractName string    `json:"contract_name"`
}

type PendingNotification struct {
	Recipient  NotificationRecipient
	Contract   *ContractNotification
	Stage      *StageNotification
	NotifyDate time.Time
}