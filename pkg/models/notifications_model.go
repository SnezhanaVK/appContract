package models

type Notification struct {
	Id_user int `json:"id_user"`

	Email                         string `json:"email"`
	Id_notification_settings      int    `json:"id_notification_settings"`
	Variant_notification_settings string `json:"variant_notification_settings"`
	Id_stage                      int    `json:"id_stage"`
	Id_contract                   int    `json:"id_contract"`
}