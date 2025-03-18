package models

type HistoreStatus struct {
	Id_history_state int    `json:"id_history_state"`
	Id_status_stage  int    `json:"id_status_stage"`
	Id_stage         int    `json:"id_stage"`
	Data_create      string `json:"data_create"`
	Comment          string `json:"comment"`
}
