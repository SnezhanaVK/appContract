package models

type Counterparty struct {
	Id_counterparty   int    `json:"id_counterparty"`
	Name_counterparty string `json:"name_counterparty"`
	Contact_info      string `json:"contact_info"`
	INN               string `json:"inn"`
	OGRN              string `json:"ogrn"`
	Adress            string `json:"adress"`
	Dop_info          string `json:"dop_info"`
}
