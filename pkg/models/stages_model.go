package models

type Stages struct {
	Id_stage          int    `json:"id_stage"`
	Name_stage        string `json:"name_stage"`
	Id_user           int    `json:"id_user"`
	Description       string `json:"description"`
	Id_status_stage   int    `json:"id_status_stage"`
	Data_create_start string `json:"data_create_start"`
	Date_create_end   string `json:"date_create_end"`
	Id_contract       int    `json:"id_contract"`

	Surname    string `json:"surname"`
	Username   string `json:"username"`
	Patronymic string `json:"patronymic"`
	Phone      string `json:"phone"`
	Photo      byte   `json:"photo"`
	Email      string `json:"email"`

	Name_status_stage string `json:"name_status_stage"`

	Name_contract        string `json:"name_contract"`
	Data_contract_create string `json:"data_contract_create"`

	Id_type_contract   int    `json:"id_type_contract"`
	Name_type_contract string `json:"d_type_contract"`

	Id_history_state int    `json:"id_history_state"`
	Data_create      string `json:"data_create"`
	Comment          string `json:"comment"`
}
