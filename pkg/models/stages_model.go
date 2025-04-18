package models

import (
	"time"
)

type Stages struct {
	Id_stage          int       `json:"id_stage"`
	Name_stage        string    `json:"name_stage"`
	Id_user           int       `json:"id_user"`
	Description       string    `json:"description"`
	Id_status_stage   int       `json:"id_status_stage"`
	Date_create_start time.Time `json:"date_create_start"`
	Date_create_end   time.Time `json:"date_create_end"`
	Id_contract       int       `json:"id_contract"`

	Surname    string `json:"surname"`
	Username   string `json:"username"`
	Patronymic string `json:"patronymic"`
	Phone      string `json:"phone"`
	Photo      byte   `json:"photo"`
	Email      string `json:"email"`

	Name_status_stage string `json:"name_status_stage"`

	Name_contract        string    `json:"name_contract"`
	Data_contract_create time.Time `json:"data_contract_create"`

	Id_type_contract   int    `json:"id_type_contract"`
	Name_type_contract string `json:"d_type_contract"`

	Id_history_state    int       `json:"id_history_state"`
	Data_create         time.Time `json:"data_create"`
	Date_change_status  time.Time `json:"date_change_status"`
	Comment             string    `json:"comment"`
	Id_comment          int       `json:"id_comment"`
	Date_create_comment time.Time `json:"date_create_comment"`
}
type StatusStage struct {
	Id_status_stage   int    `json:"id_status_stage"`
	Name_status_stage string `json:"name_status_stage"`
}

type File struct {
	Id_file   int    `json:"id_file"`
	Name_file string `json:"name_file"`
	Data      []byte `json:"data"`
	Type_file string `json:"type_file"`
	Id_stage  int    `json:"id_stage"`
}
