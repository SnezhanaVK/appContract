package db

import (
	"appContract/pkg/db"
	"appContract/pkg/models"
	"log"
)


func DBgetContractAll() ([]models.Contracts, error) {//сделать вывод информации по внешним ключам
	// соединение с бд
	conn, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// запрос к бд
	rows, err := conn.Query(`
	SELECT c.id_contract,
        c.name_contract,
        c.id_user,
        u.surname,
        u.username,
        u.patronymic,
        c.data_conclusion,
        c.data_end,
        c.id_type,
        tc.name_type_contract,
        c.id_counterparty,
        cp.name_counterparty,
        c.id_status_contract,
        sc.name_status_contract,
        cbt.id_teg,
        t.name_teg
    FROM contracts c
    JOIN 
        users u ON c.id_user = u.id_user
    JOIN 
        types_contracts tc ON c.id_type = tc.id_type_contract
    JOIN 
        counterparty cp ON c.id_counterparty = cp.id_counterparty
    JOIN 
        status_contracts sc ON c.id_status_contract = sc.id_status_contract
    JOIN 
        contracts_by_tegs cbt ON c.id_contract = cbt.id_contract
    JOIN 
        tegs t ON cbt.id_teg = t.id_teg
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// обработка результата
	var contracts []models.Contracts
	for rows.Next() {
		var contract models.Contracts	
	
		err=rows.Scan(
	&contract.Id_contract,
    &contract.Name_contract,
    &contract.User_id,
    &contract.Surname,
    &contract.Username,
    &contract.Patronymic,
    &contract.Data_conclusion,
    &contract.Data_end,
    &contract.Id_type,
    &contract.Name_type,
    &contract.Id_counterparty,
    &contract.Name_counterparty,
    &contract.Id_status_contract,
    &contract.Name_status_contract,
    &contract.Id_teg_contract,
    &contract.Tegs_contract,
		)
		if err != nil {
			log.Fatal(err)
		}
		contracts = append(contracts, contract)
	}
	return contracts, nil
}

func DBgetContractID(contractID int) ([]models.Contracts, error) {
	// соединение с бд
	conn, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// запрос к бд			
	rows, err := conn.Query(`SELECT 
		c.id_contract,
        c.name_contract,
        c.id_user,
        u.surname,
        u.username,
        u.patronymic,
        c.data_conclusion,
        c.data_end,
        c.id_type,
        tc.name_type_contract,
        c.id_counterparty,
        cp.name_counterparty,
        c.id_status_contract,
        sc.name_status_contract,
        cbt.id_teg,
        t.name_teg
    FROM 
        contracts c
    JOIN 
        users u ON c.id_user = u.id_user
    JOIN 
        types_contracts tc ON c.id_type = tc.id_type_contract
    JOIN 
        counterparty cp ON c.id_counterparty = cp.id_counterparty
    JOIN 
        status_contracts sc ON c.id_status_contract = sc.id_status_contract
    JOIN 
        contracts_by_tegs cbt ON c.id_contract = cbt.id_contract
    JOIN 
        tegs t ON cbt.id_teg = t.id_teg
		 FROM contracts WHERE id_contract=$1`,contractID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()		

	// обработка результата		
	var cotnracts []models.Contracts
	for rows.Next(){
		var contract models.Contracts
		err=rows.Scan(
	&contract.Id_contract,
    &contract.Name_contract,
    &contract.User_id,
    &contract.Surname,
    &contract.Username,
    &contract.Patronymic,
    &contract.Data_conclusion,
    &contract.Data_end,
    &contract.Id_type,
    &contract.Name_type,
    &contract.Id_counterparty,
    &contract.Name_counterparty,
    &contract.Id_status_contract,
    &contract.Name_status_contract,
    &contract.Id_teg_contract,
    &contract.Tegs_contract,
	)
		if err !=nil{
			log.Fatal(err)
		}
		cotnracts=append(cotnracts,contract)
	}
	return cotnracts, nil
}

func DBgetContractUserId(user_id int) ([]models.Contracts, error) {
    conn, err := db.ConnectDB()
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    rows, err := conn.Query(`SELECT 
            c.id_contract,
            c.name_contract,
            c.date_create_contract,
            c.id_user,
            u.surname,
            u.username,
            u.patronymic,
            u.phone,
            u.email,
            c.data_conclusion,
            c.data_end,
            c.id_type,
            tc.name_type_contract,
            c.cost,
            c.object_contract,
            c.term_payment,
            c.id_counterparty,
            cp.name_counterparty,
            cp.contact,
            cp.inn,
            cp.ogrn,
            cp.address,
            cp.dop_info,
            c.id_status_contract,
            sc.name_status_contract,
            c.notes,
            c.conditions
        FROM 
            contracts c
        JOIN 
            users u ON c.id_user = u.id_user
        JOIN 
            types_contracts tc ON c.id_type = tc.id_type_contract
        JOIN 
            counterparty cp ON c.id_counterparty = cp.id_counterparty
        JOIN 
            status_contracts sc ON c.id_status_contract = sc.id_status_contract
        WHERE c.user_id = $1`, user_id)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var contracts []models.Contracts
    for rows.Next() {
        var contract models.Contracts
        err = rows.Scan(&contract.Id_contract,
					    &contract.Name_contract, 
						&contract.Data_contract_create, 
						&contract.User_id, 
						&contract.Surname, 
						&contract.Username, 
						&contract.Patronymic, 
						&contract.Phone, 
						&contract.Email, 
						&contract.Data_conclusion, 
						&contract.Data_end, 
						&contract.Id_type, 
						&contract.Name_type, 
						&contract.Cost, 
						&contract.Object_contract, 
						&contract.Term_contract, 
						&contract.Id_counterparty, 
						&contract.Name_counterparty, 
						&contract.Contact_info, 
						&contract.Inn, 
						&contract.Ogrn, 
						&contract.Adress, 
						&contract.Dop_info, 
						&contract.Id_status_contract, 
						&contract.Name_status_contract, 
						&contract.Notes, 
						&contract.Condition)
        if err != nil {
            log.Fatal(err)
        }
        contracts = append(contracts, contract)
    }
    return contracts, nil
}

func DBaddContract(contract models.Contracts)error{
	conn ,err:=db.ConnectDB()
	if err !=nil{
		log.Fatal(err)
	}
	defer conn.Close()

	_,err=conn.Exec(`
	INSERT INTO contracts (
  name_contract,
  date_create_contract,
  id_user,
  data_conclusion,
  data_end,
  id_type,
  cost,
  object_contract,
  term_payment,
  id_counterparty,
  id_status_contract,
  notes,
  conditions
		)VALUES(
		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13
		)
		`,
		contract.Id_contract,
		contract.Name_contract,
		contract.Data_contract_create,
		contract.User_id,
		contract.Data_conclusion,
		contract.Data_end,
		contract.Id_type,
		contract.Cost,
		contract.Object_contract,
		contract.Term_contract,
		contract.Id_counterparty,
		contract.Id_status_contract,
		contract.Notes,
		contract.Condition,
	)
	if err!=nil{
		log.Fatal(err)
	}
	return nil
}

func DBchangeContract(contract models.Contracts) error{
	conn, err:= db.ConnectDB()
	if err!=nil{
		log.Fatal(err)
	}
		defer conn.Close()
	
	_, err=conn.Exec(`
	UPDATE  contracts SET
		UPDATE contracts SET

		name_contract = $1,
		date_create_contract = $2,
		id_user = $3,
		data_conclusion = $4,
		data_end = $5,
		id_type = $6,
		cost = $7,
		object_contract = $8,
		term_payment = $9,
		id_counterparty = $10,
		id_status_contract = $11,
		notes = $12,
		conditions = $13
		WHERE id_contract = $14
		`,
		contract.Id_contract,
		contract.Name_contract,
		contract.Data_contract_create,
		contract.User_id,
		contract.Data_conclusion,
		contract.Data_end,
		contract.Id_type,
		contract.Cost,
		contract.Object_contract,
		contract.Term_contract,
		contract.Id_counterparty,
		contract.Id_status_contract,
		contract.Notes,
		contract.Condition,
		contract.Id_contract,

		
	)
	if err!=nil{
		log.Fatal(err)
	}
	return nil
}

func DBchangeContractUser(id_contract int, id_user int) error{
	conn, err:= db.ConnectDB()
	if err!=nil{
		log.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Exec(`
		UPDATE contracts
		SET user_id = $2
		WHERE id_contract = $1
	`, id_contract, id_user)

	if err!=nil{
		log.Fatal(err)
	}
	return nil
}

func DBdeleteContract(contract_id int)error{
	conn, err:=db.ConnectDB()
	if err!=nil{
		log.Fatal(err)

	}
	defer conn.Close()

	_,err=conn.Exec(
		`
		DELETE FROM contracts WHERE id_contract=$1`, contract_id)
		if err!=nil{
			log.Fatal(err)
		}
		return nil
}
