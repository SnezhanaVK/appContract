package db

import (
	"appContract/pkg/db"
	"appContract/pkg/models"
	"encoding/json"
	"errors"
)

func DBgetContractAll() ([]models.Contracts, error) {
    conn:= db.GetDB()
    if conn==nil{
        return nil, errors.New("connection error")
    }

    rows, err := conn.Query(`
        SELECT 
            c.id_contract,
            c.name_contract,
            c.id_user,
            u.surname,
            u.username,
            u.patronymic,
            c.date_conclusion,
            c.date_end,
            c.date_create_contract,
            c.id_type,
            tc.name_type_contract,
            c.id_counterparty,
            cp.name_counterparty,
            c.id_status_contract,
            sc.name_status_contract,
            json_agg(json_build_object(
                'id_tegs', t.id_teg,
                'name_tegs', t.name_teg
            )) as tegs
        FROM 
            contracts c
        JOIN users u ON c.id_user = u.id_user
        JOIN types_contracts tc ON c.id_type = tc.id_type_contract
        JOIN counterparty cp ON c.id_counterparty = cp.id_counterparty
        JOIN status_contracts sc ON c.id_status_contract = sc.id_status_contract
        JOIN contracts_by_tegs cbt ON c.id_contract = cbt.id_contract
        JOIN tegs t ON cbt.id_teg = t.id_teg
        GROUP BY 
            c.id_contract, c.name_contract, c.id_user, u.surname, u.username, u.patronymic,
            c.date_conclusion, c.date_end, c.date_create_contract, c.id_type, tc.name_type_contract,
            c.id_counterparty, cp.name_counterparty, c.id_status_contract, sc.name_status_contract
        ORDER BY c.id_contract
    `)
    
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var contracts []models.Contracts
    for rows.Next() {
        var contract models.Contracts
        var tegsJSON []byte
        
        err = rows.Scan(
            &contract.Id_contract,
            &contract.Name_contract,
            &contract.Id_user,
            &contract.Surname,
            &contract.Username,
            &contract.Patronymic,
            &contract.Date_conclusion,
            &contract.Date_end,
            &contract.Date_contract_create,
            &contract.Id_type,
            &contract.Name_type,
            &contract.Id_counterparty,
            &contract.Name_counterparty,
            &contract.Id_status_contract,
            &contract.Name_status_contract,
            &tegsJSON,
        )
        
        if err != nil {
            return nil, err
        }
        
        // Декодируем JSON с тегами
        if err := json.Unmarshal(tegsJSON, &contract.Tegs); err != nil {
            return nil, err
        }
        
        contracts = append(contracts, contract)
    }
    
    return contracts, nil
}
//Sort 
func DBgetContractByType(idType int) ([]models.Contracts, error) {
    
    conn:=db.GetDB()
    if conn==nil{
        return nil, errors.New("connection error")
    }
   
    rows, err := conn.Query(`
        SELECT c.id_contract,
            c.name_contract,
            c.id_user,
            u.surname,
            u.username,
            u.patronymic,
            c.date_conclusion,
            c.date_create_contract,
            c.date_end,
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
        WHERE c.id_type = $1
    `, idType)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // обработка результата
    var contracts []models.Contracts
    for rows.Next() {
        var contract models.Contracts
        err = rows.Scan(
            &contract.Id_contract,
            &contract.Name_contract,
            &contract.Id_user,
            &contract.Surname,
            &contract.Username,
            &contract.Patronymic,
            &contract.Date_conclusion,
            &contract.Date_contract_create,
            &contract.Date_end,
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
            return nil, err
        }
        contracts = append(contracts, contract)
    }

    return contracts, nil
}
func DBgetContractsByDateCreate( date models.Date ) ([]models.Contracts, error) {
    
    conn:=db.GetDB()
   
    if conn==nil{
        return nil, errors.New("connection error")
    }
  
    rows, err := conn.Query(`
        SELECT c.id_contract,
            c.name_contract,
            c.id_user,
            u.surname,
            u.username,
            u.patronymic,
            c.date_conclusion,
            c.date_create_contract,
            c.date_end,
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
        WHERE c.date_create_contract >= $1 AND c.date_create_contract <= $2
        ORDER BY c.date_create_contract
    `, date.Date_start, date.Date_end)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // обработка результата
    var contracts []models.Contracts
    for rows.Next() {
        var contract models.Contracts
        err = rows.Scan(
            &contract.Id_contract,
            &contract.Name_contract,
            &contract.Id_user,
            &contract.Surname,
            &contract.Username,
            &contract.Patronymic,
            &contract.Date_conclusion,
            &contract.Date_contract_create,
            &contract.Date_end,
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
            return nil, err
        }
        contracts = append(contracts, contract)
    }

    return contracts, nil
}
func DBgetContractsByTegs() ([]models.Contracts, error) {
   
	conn := db.GetDB()

    if conn == nil {
        return nil, errors.New("DB connection is nil")
    }
	
	rows, err := conn.Query(`
	   SELECT c.id_contract,
            c.name_contract,
            c.id_user,
            u.surname,
            u.username,
            u.patronymic,
            c.date_conclusion,
            c.date_create_contract,
            c.date_end,
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
        ORDER BY cbt.id_teg
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// обработка результата
	var contracts []models.Contracts
	for rows.Next() {
		var contract models.Contracts	
		err=rows.Scan(
		&contract.Id_contract,
		&contract.Name_contract,
		&contract.Id_user,
		&contract.Surname,
		&contract.Username,
		&contract.Patronymic,
		&contract.Date_conclusion,
		&contract.Date_contract_create,
		&contract.Date_end,
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
			return nil, err
		}
		contracts = append(contracts, contract)
	}
	
	return contracts, nil
}

func DBgetContractsByStatus() ([]models.Contracts, error) {
    
    conn:= db.GetDB()
    if conn==nil{
        return nil, errors.New("connection error")
    }
    
    rows, err := conn.Query(`
        SELECT c.id_contract,
            c.name_contract,
            c.id_user,
            u.surname,
            u.username,
            u.patronymic,
            c.date_conclusion,
            c.date_create_contract,
            c.date_end,
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
        ORDER BY c.id_status_contract
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // обработка результата
    var contracts []models.Contracts
    for rows.Next() {
        var contract models.Contracts
        err = rows.Scan(
            &contract.Id_contract,
            &contract.Name_contract,
            &contract.Id_user,
            &contract.Surname,
            &contract.Username,
            &contract.Patronymic,
            &contract.Date_conclusion,
            &contract.Date_contract_create,
            &contract.Date_end,
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
            return nil, err
        }
        contracts = append(contracts, contract)
    }

    return contracts, nil
}


func DBgetContractID(contractID int) ([]models.Contracts, error) {
    conn:= db.GetDB()
    if conn==nil{
        return nil, errors.New("connection error")
    }
    rows, err := conn.Query(`
        SELECT 
            c.id_contract,
            c.name_contract,
            c.id_user,
            u.surname,
            u.username,
            u.patronymic,
            c.date_conclusion,
            c.date_end,
            c.date_create_contract,
            c.id_type,
            tc.name_type_contract,
            c.id_counterparty,
            cp.name_counterparty,
            c.id_status_contract,
            sc.name_status_contract,
            json_agg(json_build_object(
                'id_tegs', t.id_teg,  
                'name_tegs', t.name_teg  
            )) as tegs
        FROM 
            contracts c
        JOIN users u ON c.id_user = u.id_user
        JOIN types_contracts tc ON c.id_type = tc.id_type_contract
        JOIN counterparty cp ON c.id_counterparty = cp.id_counterparty
        JOIN status_contracts sc ON c.id_status_contract = sc.id_status_contract
        JOIN contracts_by_tegs cbt ON c.id_contract = cbt.id_contract
        JOIN tegs t ON cbt.id_teg = t.id_teg
        WHERE c.id_contract = $1
        GROUP BY 
            c.id_contract, c.name_contract, c.id_user, u.surname, u.username, u.patronymic,
            c.date_conclusion, c.date_end, c.date_create_contract, c.id_type, tc.name_type_contract,
            c.id_counterparty, cp.name_counterparty, c.id_status_contract, sc.name_status_contract`,
        contractID)
    
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var contracts []models.Contracts
    for rows.Next() {
        var contract models.Contracts
        var tegsJSON []byte
        
        err = rows.Scan(
            &contract.Id_contract,
            &contract.Name_contract,
            &contract.Id_user,
            &contract.Surname,
            &contract.Username,
            &contract.Patronymic,
            &contract.Date_conclusion,
            &contract.Date_end,
            &contract.Date_contract_create,
            &contract.Id_type,
            &contract.Name_type,
            &contract.Id_counterparty,
            &contract.Name_counterparty,
            &contract.Id_status_contract,
            &contract.Name_status_contract,
            &tegsJSON,
        )
        
        if err != nil {
            return nil, err
        }
        
        // Декодируем JSON с тегами
        if err := json.Unmarshal(tegsJSON, &contract.Tegs); err != nil {
            return nil, err
        }
        
        contracts = append(contracts, contract)
    }
    
    return contracts, nil
}
func DBgetContractUserId(user_id int) ([]models.Contracts, error) {
    conn:= db.GetDB()

    if conn==nil{
        return nil, errors.New("connection error")
    }
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
            c.date_conclusion,
            c.date_end,
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
        WHERE c.id_user = $1`, user_id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var contracts []models.Contracts
    for rows.Next() {
        var contract models.Contracts
        err = rows.Scan(&contract.Id_contract,
					    &contract.Name_contract, 
						&contract.Date_contract_create, 
						&contract.Id_user, 
						&contract.Surname, 
						&contract.Username, 
						&contract.Patronymic, 
						&contract.Phone, 
						&contract.Email, 
						&contract.Date_conclusion, 
						&contract.Date_end, 
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
            return nil, err
        }
        contracts = append(contracts, contract)
    }
    return contracts, nil
}

func DBaddContract(contract models.Contracts)error{
	conn:= db.GetDB()
	if conn==nil{
        return errors.New("connection error")
    }
	var userExist bool
    err := conn.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id_user = $1)`, contract.Id_user).Scan(&userExist)
    if err != nil {
        return err
    }
    if !userExist {
        return errors.New("user not found")
    }

	_,err=conn.Exec(`
	INSERT INTO contracts (
  name_contract,
  date_create_contract,
  id_user,
  date_conclusion,
  date_end,
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
		contract.Name_contract,
		contract.Date_contract_create,
		contract.Id_user,
		contract.Date_conclusion,
		contract.Date_end,
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
		
		return err
	}
	return nil
}

func DBchangeContract(contract models.Contracts) error{
	conn:= db.GetDB()
	if conn==nil{
        return errors.New("connection error")
    }

	_, err:=conn.Exec(`
	UPDATE  contracts SET
		name_contract = $1,
		date_create_contract = $2,
		id_user = $3,
		date_conclusion = $4,
		date_end = $5,
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
		
		contract.Name_contract,
		contract.Date_contract_create,
		contract.Id_user,
		contract.Date_conclusion,
		contract.Date_end,
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
		return err
	}
	return nil
}

func DBchangeContractUser(id_contract int, id_user int) error {
    conn:= db.GetDB()
    if conn==nil{
        return errors.New("connection error")
    }
    
    var exists bool
    err := conn.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id_user = $1)`, id_user).Scan(&exists)
    if err != nil {
        return err
    }
    if !exists {
        return errors.New("id_user не существует в таблице users")
    }

   
    exists = false
    err = conn.QueryRow(`SELECT EXISTS(SELECT 1 FROM contracts WHERE id_contract = $1)`, id_contract).Scan(&exists)
    if err != nil {
        return err
    }
    if !exists {
        return errors.New("id_contract не существует в таблице contracts")
    }


    result, err := conn.Exec(`
        UPDATE contracts
        SET id_user = $2
        WHERE id_contract = $1
    `, id_contract, id_user)

    if err != nil {
        return err
    }

    rowsAffected := result.RowsAffected()
    if rowsAffected == 0 {
        return errors.New("id_contract или id_user не существует в таблице contracts")
    }

    return nil
}

func DBdeleteContract(contract_id int)error{
	conn:=db.GetDB()
    if conn==nil{
        return errors.New("connection error")
    }
	
	_,err:=conn.Exec(
		`
		DELETE FROM contracts WHERE id_contract=$1`, contract_id)
		if err!=nil{
			return err
		}
		return nil
}
