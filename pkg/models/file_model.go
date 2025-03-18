package models

type File struct {
	Id_file   int    `json:"id_file"`
	Name_file string `json:"name_file"`
	Data      byte   `json:"data"`
	Type_file string `json:"type_file"`
	Id_stage  int    `json:"id_stage"`
}