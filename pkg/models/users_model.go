package models

type Users struct {
	Id_user    int    `json:"id_user"`
	Surname    string `json:"surname"`
	Username   string `json:"username"`
	Patronymic string `json:"patronymic"`
	Phone      string `json:"phone"`
	Photo      string `json:"photo"`
	Email      string `json:"email"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	Admin      bool   `json:"admin"`
	Meneger    bool   `json:"meneger"`
	Roles      []Role `json:"roles"`

	//Notification
	Id_notification_settings      int    `json:"id_notification"`
	Variant_notification_settings string `json:"name_notification"`

	//Role
	Id_role   int    `json:"id_role"`
	Name_role string `json:"name_role"`
}

type Role struct {
	Id_role   int    `json:"id_role"`
	Name_role string `json:"name_role"`
}
