package models

import (
	"time"
)

type Contracts struct {
	
	Id_contract     int       `json:"id_contract"`
	Name_contract   string    `json:"name_contract"`
	Date_contract_create   time.Time`json:"date_contract_create"`
	Id_user         int       `json:"id_user"`
	Date_conclusion time.Time `json:"date_conclusion"`
	Date_end        time.Time `json:"date_end"`
	Id_type         int       `json:"id_type"`
	Cost            int       `json:"cost"`
	Object_contract string    `json:"object_contract"`
	Term_contract   string    `json:"term_contract"`
	Id_counterparty int       `json:"id_counterparty"`
	Id_status_contract int     `json:"id_status_contract"`
	Notes           string    `json:"notes"`
	Condition       string    `json:"condition"`
	Tegs [] Tag 			  `json:"tegs"`

	Surname     string        `json:"surname"`
    Username        string    `json:"username"`
    Patronymic      string    `json:"patronymic"`
	Name_type       string    `json:"name_type"`
	Name_counterparty string  `json:"name_counterparty"`
	Name_status_contract string`json:"name_status_contract"`
	Tegs_contract     string  `json:"tegs_contract"`
	Phone           string    `json:"phone"`
	Email           string    `json:"email"`
	Contact_info    string    `json:"contact_info"`
	Inn             string    `json:"inn"`
	Ogrn            string    `json:"ogrn"`
	Adress          string    `json:"adress"`
	Dop_info        string    `json:"dop_info"`
	Id_teg_contract int       `json:"id_teg_contract"`

}

type Tag struct {
	Id_tegs   int    `json:"id_tegs"`
	Name_tegs string `json:"name_tegs"`
}

type Date struct {
    Date_start time.Time `json:"date_start"`

    Date_end time.Time `json:"date_end"`
}

// type Type_contracts struct {
// 	Id_type_contract    int    `json:"id_type_contract"`
// 	Name_type_contract   string `json:"d_type_contract"`
	
// }

// type Status_contracts struct {
// 	Id_status_contract int    `json:"id_status_contract"`
// 	Name_status_contract        string `json:"name_status_contract"`
// }