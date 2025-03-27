package models

type Users struct {
	Id_user         int    `json:"id_user"`
	Surname         string `json:"surname"`
	Username        string `json:"username"`
	Patronymic      string `json:"patronymic"`
	Phone           string `json:"phone"`
	Photo           []byte `json:"photo"`
	Email           string `json:"email"`
	Role_id         int    `json:"role_id"`
	Notification_id int    `json:"notification_id"`
	Admin           bool   `json:"admin"`
	Login           string `json:"login"`
	Password        string `json:"password"`

	Variant_notification string `json:"variant_notification"`
}