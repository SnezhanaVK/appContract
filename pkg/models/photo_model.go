package models

type Photo struct {
	Id_photo   int    `json:"id_photo"`
	Data_photo []byte `json:"date_photo"`
	Type_photo string `json:"type_photo"`
	Id_user    int    `json:"id_user"`
}
