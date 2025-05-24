package db

//contracts_db.go

import (
	"appContract/pkg/db"
	"appContract/pkg/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

func DBgetContractAll() ([]models.Contracts, error) {
	conn := db.GetDB()
	if conn == nil {
		return nil, errors.New("connection error")
	}

	rows, err := conn.Query(context.Background(), `
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
    COALESCE(
        json_agg(
            json_build_object(
                'id_tegs', t.id_teg,
                'name_tegs', t.name_teg
            )
        ) FILTER (WHERE t.id_teg IS NOT NULL),
        '[]'::json
    ) as tegs
FROM 
    contracts c
JOIN users u ON c.id_user = u.id_user
JOIN types_contracts tc ON c.id_type = tc.id_type_contract
JOIN counterparty cp ON c.id_counterparty = cp.id_counterparty
JOIN status_contracts sc ON c.id_status_contract = sc.id_status_contract
LEFT JOIN contracts_by_tegs cbt ON c.id_contract = cbt.id_contract
LEFT JOIN tegs t ON cbt.id_teg = t.id_teg
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

		if err := json.Unmarshal(tegsJSON, &contract.Tags); err != nil {
			return nil, err
		}

		contracts = append(contracts, contract)
	}

	return contracts, nil
}

// Sort
func DBgetContractByType(idType int) ([]models.Contracts, error) {
	conn := db.GetDB()
	if conn == nil {
		return nil, errors.New("connection error")
	}
	rows, err := conn.Query(context.Background(), `
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
            COALESCE(
                json_agg(
                    json_build_object(
                        'id_tegs', t.id_teg,
                        'name_tegs', t.name_teg
                    )
                ) FILTER (WHERE t.id_teg IS NOT NULL),
                '[]'::json
            ) as tegs
        FROM contracts c
        JOIN users u ON c.id_user = u.id_user
        JOIN types_contracts tc ON c.id_type = tc.id_type_contract
        JOIN counterparty cp ON c.id_counterparty = cp.id_counterparty
        JOIN status_contracts sc ON c.id_status_contract = sc.id_status_contract
        LEFT JOIN contracts_by_tegs cbt ON c.id_contract = cbt.id_contract
        LEFT JOIN tegs t ON cbt.id_teg = t.id_teg
        WHERE c.id_type = $1
        GROUP BY 
            c.id_contract, c.name_contract, c.id_user, u.surname, u.username, u.patronymic,
            c.date_conclusion, c.date_end, c.date_create_contract, c.id_type, tc.name_type_contract,
            c.id_counterparty, cp.name_counterparty, c.id_status_contract, sc.name_status_contract
        ORDER BY c.id_contract
    `, idType)

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

		if err := json.Unmarshal(tegsJSON, &contract.Tags); err != nil {
			return nil, err
		}

		contracts = append(contracts, contract)
	}

	return contracts, nil
}
func DBgetContractsByDateCreate(date models.Date) ([]models.Contracts, error) {
	conn := db.GetDB()
	if conn == nil {
		return nil, errors.New("connection error")
	}

	rows, err := conn.Query(context.Background(), `
        SELECT 
            c.id_contract,
            c.name_contract,
            c.id_user,
            u.surname,
            u.username,
            u.patronymic,
            c.date_conclusion,
            c.date_create_contract,
            c.id_type,
            tc.name_type_contract,
            c.id_counterparty,
            cp.name_counterparty,
            c.id_status_contract,
            sc.name_status_contract,
            COALESCE(
                json_agg(
                    json_build_object(
                        'id_tegs', t.id_teg,
                        'name_tegs', t.name_teg
                    )
                ) FILTER (WHERE t.id_teg IS NOT NULL),
                '[]'::json
            ) as tegs
        FROM contracts c
        JOIN users u ON c.id_user = u.id_user
        JOIN types_contracts tc ON c.id_type = tc.id_type_contract
        JOIN counterparty cp ON c.id_counterparty = cp.id_counterparty
        JOIN status_contracts sc ON c.id_status_contract = sc.id_status_contract
        LEFT JOIN contracts_by_tegs cbt ON c.id_contract = cbt.id_contract
        LEFT JOIN tegs t ON cbt.id_teg = t.id_teg
        WHERE c.date_create_contract BETWEEN $1 AND $2
        GROUP BY 
            c.id_contract, c.name_contract, c.id_user, u.surname, u.username, u.patronymic,
            c.date_conclusion, c.date_create_contract, c.date_end, c.id_type, tc.name_type_contract,
            c.id_counterparty, cp.name_counterparty, c.id_status_contract, sc.name_status_contract
        ORDER BY c.date_create_contract
    `, date.Date_start, date.Date_end)

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

		if err := json.Unmarshal(tegsJSON, &contract.Tags); err != nil {
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

	rows, err := conn.Query(context.Background(), `
        SELECT 
            c.id_contract,
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
            COALESCE(
                json_agg(
                    json_build_object(
                        'id_tegs', t.id_teg,
                        'name_tegs', t.name_teg
                    )
                ) FILTER (WHERE t.id_teg IS NOT NULL),
                '[]'::json
            ) as tegs
        FROM contracts c
        JOIN users u ON c.id_user = u.id_user
        JOIN types_contracts tc ON c.id_type = tc.id_type_contract
        JOIN counterparty cp ON c.id_counterparty = cp.id_counterparty
        JOIN status_contracts sc ON c.id_status_contract = sc.id_status_contract
        LEFT JOIN contracts_by_tegs cbt ON c.id_contract = cbt.id_contract
        LEFT JOIN tegs t ON cbt.id_teg = t.id_teg
        GROUP BY 
            c.id_contract, c.name_contract, c.id_user, u.surname, u.username, u.patronymic,
            c.date_conclusion, c.date_create_contract, c.date_end, c.id_type, tc.name_type_contract,
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
			&contract.Date_contract_create,
			&contract.Date_end,
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

		if err := json.Unmarshal(tegsJSON, &contract.Tags); err != nil {
			return nil, err
		}

		contracts = append(contracts, contract)
	}

	return contracts, nil
}

func DBgetContractsByStatus() ([]models.Contracts, error) {
	conn := db.GetDB()
	if conn == nil {
		return nil, errors.New("connection error")
	}

	rows, err := conn.Query(context.Background(), `
        SELECT 
            c.id_contract,
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
            COALESCE(
                json_agg(
                    json_build_object(
                        'id_tegs', t.id_teg,
                        'name_tegs', t.name_teg
                    )
                ) FILTER (WHERE t.id_teg IS NOT NULL),
                '[]'::json
            ) as tegs
        FROM contracts c
        JOIN users u ON c.id_user = u.id_user
        JOIN types_contracts tc ON c.id_type = tc.id_type_contract
        JOIN counterparty cp ON c.id_counterparty = cp.id_counterparty
        JOIN status_contracts sc ON c.id_status_contract = sc.id_status_contract
        LEFT JOIN contracts_by_tegs cbt ON c.id_contract = cbt.id_contract
        LEFT JOIN tegs t ON cbt.id_teg = t.id_teg
        GROUP BY 
            c.id_contract, c.name_contract, c.id_user, u.surname, u.username, u.patronymic,
            c.date_conclusion, c.date_create_contract, c.date_end, c.id_type, tc.name_type_contract,
            c.id_counterparty, cp.name_counterparty, c.id_status_contract, sc.name_status_contract
        ORDER BY c.id_status_contract
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
			&contract.Date_contract_create,
			&contract.Date_end,
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

		if err := json.Unmarshal(tegsJSON, &contract.Tags); err != nil {
			return nil, err
		}

		contracts = append(contracts, contract)
	}

	return contracts, nil
}

func DBgetContractID(contractID int) ([]models.Contracts, error) {
	conn := db.GetDB()
	if conn == nil {
		return nil, errors.New("connection error")
	}
	rows, err := conn.Query(context.Background(), `
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
			c.notes,
			c.conditions,
			c.cost,
			c.object_contract,
			c.term_payment,
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
        Left JOIN contracts_by_tegs cbt ON c.id_contract = cbt.id_contract
        Left JOIN tegs t ON cbt.id_teg = t.id_teg
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
			&contract.Notes,      
			&contract.Conditions, 
			&contract.Cost,      
			&contract.Object_contract,
			&contract.Term_payment,
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

		if err := json.Unmarshal(tegsJSON, &contract.Tags); err != nil {
			return nil, err
		}

		contracts = append(contracts, contract)
	}

	return contracts, nil
}

func DBgetContractUserId(user_id int) ([]models.Contracts, error) {
	conn := db.GetDB()
	if conn == nil {
		return nil, errors.New("connection error")
	}

	query := `
        SELECT 
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
            c.conditions,
            COALESCE(
                json_agg(
                    json_build_object(
                        'id_tegs', t.id_teg,
                        'name_tegs', t.name_teg
                    )
                ) FILTER (WHERE t.id_teg IS NOT NULL),
                '[]'::json
            ) as tegs
        FROM 
            contracts c
        JOIN users u ON c.id_user = u.id_user
        JOIN types_contracts tc ON c.id_type = tc.id_type_contract
        JOIN counterparty cp ON c.id_counterparty = cp.id_counterparty
        JOIN status_contracts sc ON c.id_status_contract = sc.id_status_contract
        LEFT JOIN contracts_by_tegs cbt ON c.id_contract = cbt.id_contract
        LEFT JOIN tegs t ON cbt.id_teg = t.id_teg
        WHERE c.id_user = $1
        GROUP BY 
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
            c.conditions`

	rows, err := conn.Query(context.Background(), query, user_id)
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
			&contract.Term_payment,
			&contract.Id_counterparty,
			&contract.Name_counterparty,
			&contract.Contact_info,
			&contract.Inn,
			&contract.Ogrn,
			&contract.Address,
			&contract.Dop_info,
			&contract.Id_status_contract,
			&contract.Name_status_contract,
			&contract.Notes,
			&contract.Conditions,
			&tegsJSON,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(tegsJSON, &contract.Tags); err != nil {
			return nil, err
		}

		contracts = append(contracts, contract)
	}

	return contracts, nil
}

func DBaddContract(contract models.Contracts) (int, error) {
	conn := db.GetDB()
	if conn == nil {
		return 0, errors.New("connection error")
	}

	// Проверка существования пользователя
	var userExist bool
	err := conn.QueryRow(context.Background(), `SELECT EXISTS(SELECT 1 FROM users WHERE id_user = $1)`, contract.Id_user).Scan(&userExist)
	if err != nil {
		return 0, err
	}
	if !userExist {
		return 0, errors.New("user not found")
	}

	// Добавляем RETURNING id_contract
	var id int
	err = conn.QueryRow(context.Background(), `
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
    ) VALUES (
        $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13
    )
    RETURNING id_contract
    `,
		contract.Name_contract,
		contract.Date_contract_create,
		contract.Id_user,
		contract.Date_conclusion,
		contract.Date_end,
		contract.Id_type,
		contract.Cost,
		contract.Object_contract,
		contract.Term_payment,
		contract.Id_counterparty,
		contract.Id_status_contract,
		contract.Notes,
		contract.Conditions,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func DBchangeContract(contract models.Contracts) error {
    conn := db.GetDB()
    if conn == nil {
        return errors.New("connection error")
    }

    _, err := conn.Exec(context.Background(), `
        UPDATE contracts SET
            name_contract = $1,
            date_conclusion = $2,
            date_end = $3,
            id_type = $4,
            cost = $5,
            object_contract = $6,
            term_payment = $7,
            id_counterparty = $8,
            id_status_contract = $9,
            notes = $10,
            conditions = $11
        WHERE id_contract = $12
    `,
        contract.Name_contract,
        contract.Date_conclusion,
        contract.Date_end,
        contract.Id_type,
        contract.Cost,
        contract.Object_contract,
        contract.Term_payment,
        contract.Id_counterparty,
        contract.Id_status_contract,
        contract.Notes,
        contract.Conditions,
        contract.Id_contract,
    )

    if err != nil {
        return fmt.Errorf("failed to update contract: %v", err)
    }
    return nil
}
func DBchangeContractUser(id_contract int, id_user int) error {
	ctx := context.Background()
	log.Printf("[DEBUG] DBchangeContractUser: contract=%d, user=%d", id_contract, id_user)

	conn := db.GetDB()
	if conn == nil {
		log.Println("[ERROR] DB connection is nil")
		return errors.New("database connection failed")
	}

	var userExists bool
	err := conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id_user = $1)`, id_user).Scan(&userExists)
	if err != nil {
		log.Printf("[ERROR] User check query failed: %v", err)
		return fmt.Errorf("user verification failed")
	}

	if !userExists {
		log.Printf("[WARN] User %d not found", id_user)
		return fmt.Errorf("user %d does not exist", id_user)
	}
	var contractExists bool
	err = conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM contracts WHERE id_contract = $1)`, id_contract).Scan(&contractExists)
	if err != nil {
		log.Printf("[ERROR] Contract check query failed: %v", err)
		return fmt.Errorf("contract verification failed")
	}

	if !contractExists {
		log.Printf("[WARN] Contract %d not found", id_contract)
		return fmt.Errorf("contract %d does not exist", id_contract)
	}
	tag, err := conn.Exec(ctx,
		`UPDATE contracts SET id_user = $2 WHERE id_contract = $1`,
		id_contract, id_user)
	if err != nil {
		log.Printf("[ERROR] Update failed: %v", err)
		return fmt.Errorf("update operation failed")
	}

	if rowsAffected := tag.RowsAffected(); rowsAffected == 0 {
		log.Printf("[WARN] No rows affected for contract %d", id_contract)
		return fmt.Errorf("no changes made to contract %d", id_contract)
	}

	log.Printf("[INFO] Successfully updated contract %d with user %d", id_contract, id_user)
	return nil
}

func DBdeleteContract(contract_id int) error {
	conn := db.GetDB()
	if conn == nil {
		return errors.New("connection error")
	}
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `
        DELETE FROM files 
        WHERE id_stage IN (
            SELECT id_stage FROM stages WHERE id_contract = $1
        )`, contract_id)
	if err != nil {
		return fmt.Errorf("error deleting files: %v", err)
	}

	_, err = tx.Exec(context.Background(), `
        DELETE FROM comments 
        WHERE id_history_status IN (
            SELECT id_history_status FROM history_status 
            WHERE id_stage IN (
                SELECT id_stage FROM stages WHERE id_contract = $1
            )
        )`, contract_id)
	if err != nil {
		return fmt.Errorf("error deleting comments: %v", err)
	}

	_, err = tx.Exec(context.Background(), `
        DELETE FROM history_status 
        WHERE id_stage IN (
            SELECT id_stage FROM stages WHERE id_contract = $1
        )`, contract_id)
	if err != nil {
		return fmt.Errorf("error deleting history_status: %v", err)
	}

	_, err = tx.Exec(context.Background(), `
        DELETE FROM stages 
        WHERE id_contract = $1`, contract_id)
	if err != nil {
		return fmt.Errorf("error deleting stages: %v", err)
	}

	_, err = tx.Exec(context.Background(), `
        DELETE FROM contracts_by_tegs 
        WHERE id_contract = $1`, contract_id)
	if err != nil {
		return fmt.Errorf("error deleting contracts_by_tegs: %v", err)
	}

	_, err = tx.Exec(context.Background(), `
        DELETE FROM contracts 
        WHERE id_contract = $1`, contract_id)
	if err != nil {
		return fmt.Errorf("error deleting contract: %v", err)
	}

	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
