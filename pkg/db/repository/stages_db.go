package db

import (
	"appContract/pkg/db"
	"appContract/pkg/models"
	"context"
	"errors"
	"fmt"
	"log"
)

func DBgetStageAll() ([]models.Stages, error) {
	conn:= db.GetDB()
	if conn == nil {
		return nil, errors.New("DB connection is nil")
	}
    rows, err := conn.Query( context.Background(),`SELECT 
    s.id_stage,
    s.name_stage,
    s.id_user,
    u.surname,
    u.username,
    u.patronymic,
    u.phone,
    u.email,
    s.description,
    hs.id_status_stage,
    ss.name_status_stage,
    s.date_create_start,
    s.date_create_end,
    s.id_contract,
    c.name_contract,
	c.date_create_contract
FROM stages s
JOIN history_status hs ON s.id_stage = hs.id_stage
JOIN users u ON s.id_user = u.id_user
JOIN contracts c ON s.id_contract = c.id_contract
JOIN status_stages ss ON hs.id_status_stage = ss.id_status_stage`)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var stages []models.Stages
	for rows.Next() {
		var stage models.Stages
		err = rows.Scan(&stage.Id_stage,
			&stage.Name_stage,
			&stage.Id_user,
			&stage.Surname,
			&stage.Username,
			&stage.Patronymic,
			&stage.Phone,
			&stage.Email,
			&stage.Description,
			&stage.Id_status_stage,
			&stage.Name_status_stage,
			&stage.Date_create_start,
			&stage.Date_create_end,
			&stage.Id_contract,
			&stage.Name_contract,
			&stage.Data_contract_create)
		if err != nil {
			log.Fatal(err)
		}
		stages = append(stages, stage)
	}
	return stages, nil
}

func DBgetStageByContractID(id_contract int) ([]models.Stages, error) {
    conn := db.GetDB()
    if conn == nil {
        return nil, errors.New("DB connection is nil")
    }

    rows, err := conn.Query( context.Background(),`
        SELECT 
            s.id_stage,
            s.name_stage,
            s.id_user,
            u.surname,
            u.username,
            u.patronymic,
            s.description,
            s.date_create_start,
            s.date_create_end,
            s.id_contract,
            c.name_contract,
            c.date_create_contract,
            latest_status.id_status_stage,
            ss.name_status_stage,
            latest_status.data_change_status
        FROM stages s
        JOIN contracts c ON s.id_contract = c.id_contract
        JOIN users u ON s.id_user = u.id_user
        JOIN (
            SELECT 
                hs.id_stage,
                hs.id_status_stage,
                hs.data_change_status,
                ROW_NUMBER() OVER (PARTITION BY hs.id_stage ORDER BY hs.data_change_status DESC) as rn
            FROM history_status hs
        ) latest_status ON s.id_stage = latest_status.id_stage AND latest_status.rn = 1
        JOIN status_stages ss ON latest_status.id_status_stage = ss.id_status_stage
        WHERE s.id_contract = $1
        ORDER BY s.id_stage
    `, id_contract)
    
    if err != nil {
        return nil, fmt.Errorf("query failed: %w", err)
    }
    defer rows.Close()

    var stages []models.Stages
    for rows.Next() {
        var stage models.Stages
        err := rows.Scan(
            &stage.Id_stage,
            &stage.Name_stage,
            &stage.Id_user,
            &stage.Surname,
            &stage.Username,
            &stage.Patronymic,
            &stage.Description,
            &stage.Data_create,
            &stage.Date_create_end,
            &stage.Id_contract,
            &stage.Name_contract,
            &stage.Data_contract_create,
            &stage.Id_status_stage,
            &stage.Name_status_stage,
            &stage.Date_change_status,
        )
        if err != nil {
            return nil, fmt.Errorf("scan failed: %w", err)
        }
        stages = append(stages, stage)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %w", err)
    }

    return stages, nil
}

// db.go - Измененный SQL-запрос
func DBgetStageUserID(user_id int) ([]models.Stages, error) {
    conn := db.GetDB()
    if conn == nil {
        return nil, errors.New("database connection is nil")
    }

    rows, err := conn.Query(context.Background(), `
        SELECT 
            s.id_stage,
            s.name_stage,
            s.description,
            s.date_create_start,
            s.date_create_end,
            c.name_contract,
            ss.name_status_stage,
            u.surname,
            u.username,
            u.patronymic
        FROM stages s
        JOIN contracts c ON s.id_contract = c.id_contract
        JOIN (
            SELECT 
                id_stage, 
                MAX(id_history_status) as last_status
            FROM history_status
            GROUP BY id_stage
        ) last_hs ON s.id_stage = last_hs.id_stage
        JOIN history_status hs ON last_hs.last_status = hs.id_history_status
        JOIN status_stages ss ON hs.id_status_stage = ss.id_status_stage
        JOIN users u ON s.id_user = u.id_user
        WHERE s.id_user = $1
    `, user_id)

    if err != nil {
        return nil, fmt.Errorf("query error: %w", err)
    }
    defer rows.Close()

    var stages []models.Stages
    for rows.Next() {
        var stage models.Stages
        if err := rows.Scan(
            &stage.Id_stage,
            &stage.Name_stage,
            &stage.Description,
            &stage.Date_create_start,
            &stage.Date_create_end,
            &stage.Name_contract,
            &stage.Name_status_stage,
            &stage.Surname,
            &stage.Username,
            &stage.Patronymic,
        ); err != nil {
            return nil, fmt.Errorf("scan error: %w", err)
        }
        stages = append(stages, stage)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %w", err)
    }

    return stages, nil
}

func DBgetStageID(stage_id int) (models.Stages, error) {
	conn:= db.GetDB()
	if conn == nil {
		return models.Stages{}, errors.New("DB connection is nil")
	}

	rows, err := conn.Query( context.Background(),`SELECT 
    s.id_stage,
    s.name_stage,
    s.id_user,
    u.surname,
    u.username,
    u.patronymic,
    u.phone,
    u.email,
    s.description,
    hs.id_status_stage,
    ss.name_status_stage,
    s.date_create_start,
    s.date_create_end,
    s.id_contract,
    c.name_contract,
    c.date_create_contract,
    t.id_type_contract,
    t.name_type_contract
FROM stages s
JOIN history_status hs ON s.id_stage = hs.id_stage
JOIN users u ON s.id_user = u.id_user
JOIN contracts c ON s.id_contract = c.id_contract
JOIN types_contracts t ON c.id_type = t.id_type_contract
JOIN status_stages ss ON hs.id_status_stage = ss.id_status_stage 
WHERE s.id_stage=$1`, stage_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var stage models.Stages
	for rows.Next() {
		err = rows.Scan(&stage.Id_stage,
			&stage.Name_stage,
			&stage.Id_user,
			&stage.Surname,
			&stage.Username,
			&stage.Patronymic,
			&stage.Phone,
			&stage.Email,
			&stage.Description,
			&stage.Id_status_stage,
			&stage.Name_status_stage,
			&stage.Date_create_start,
			&stage.Date_create_end,
			&stage.Id_contract,
			&stage.Name_contract,
			&stage.Data_contract_create,
			&stage.Id_type_contract,
			&stage.Name_type_contract)
		if err != nil {
			log.Fatal(err)
		}
	}
	return stage, nil
}
func DBgetFileIDStageID(id_stages int, id_file int) (models.File, error) {
	conn:= db.GetDB()	
	if conn == nil {
		return models.File{}, errors.New("DB connection is nil")
	}

	rows, err := conn.Query( context.Background(),"SELECT * FROM files WHERE Id_stage=$1 AND Id_file=$2", id_stages, id_file)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

    var data []byte
	var file models.File
	for rows.Next() {
		err = rows.Scan(
            &file.Id_file,
			&file.Name_file,
			&data,
			&file.Type_file,
			&file.Id_stage)
        file.Data = data
		if err != nil {
			log.Fatal(err)
		}
        
	}
	return file, nil
}

func DBgetFilesStageID(id_stages int) ([]models.File, error) {
	conn := db.GetDB()
	if conn == nil {
		return nil, errors.New("DB connection is nil")
	}

	rows, err := conn.Query( context.Background(),"SELECT * FROM files WHERE Id_stage=$1", id_stages)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

    var data []byte
	var files []models.File
	for rows.Next() {
        var file models.File
		err = rows.Scan(
            &file.Id_file,
			&file.Name_file,
			&data,
			&file.Type_file,
			&file.Id_stage)
        file.Data = data
		if err != nil {
			log.Fatal(err)
		}
        files = append(files, file)
	}
	return files, nil
}

func DBgetStageIdStatus(id_stage int) (models.StatusStage, error) {
	conn:= db.GetDB()
	if conn == nil {
		return models.StatusStage{}, errors.New("DB connection is nil")
	}
	

	rows, err := conn.Query( context.Background(),`SELECT * 
                            FROM status_stages
                            WHERE id_status_stage=$1`, id_stage)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var status models.StatusStage
	for rows.Next() {
		err = rows.Scan(&status.Id_status_stage,
			&status.Name_status_stage)
		if err != nil {
			log.Fatal(err)
		}
	}
	return status, nil
}

func DBaddFile(file models.File) error {
	conn:= db.GetDB()
	if conn == nil {
		return errors.New("DB connection is nil")
	}
	

	_, err := conn.Exec( context.Background(),`INSERT INTO files (
        name_file,
        data,
        type_file,
        id_stage
    ) VALUES ($1,$2,$3,$4)`,
    file.Name_file,
    file.Data,
    file.Type_file,
    file.Id_stage)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}
func DBaddStage(stage models.Stages) error {
	conn:= db.GetDB()
	if conn == nil {
		return errors.New("DB connection is nil")
	}

	_, err := conn.Exec( context.Background(),`INSERT INTO stages(
    name_stage,
    id_user,
    description,
    date_create_start,
    date_create_end,
    id_contract
    )VALUES ($1,$2,$3,$4,$5,$6)`,
		stage.Name_stage,
		stage.Id_user,
		stage.Description,
		stage.Date_create_start,
		stage.Date_create_end,
		stage.Id_contract)

	if err != nil {
		log.Fatal(err)
	}
	return nil

}

func DBaddComment(idStage int, idStatusStage int, comment string) error {
	conn:= db.GetDB()
	if conn == nil {
		return errors.New("DB connection is nil")
	}

    var idHistoryState int
    err := conn.QueryRow( context.Background(),`SELECT id_history_status FROM history_status
        WHERE id_stage = $1 AND id_status_stage = $2`,
        idStage,
        idStatusStage).Scan(&idHistoryState)

    if err != nil {
        log.Fatal(err)
    }

    _, err = conn.Exec( context.Background(),`INSERT INTO comments (id_history_state, comment, date_create_comment)
        VALUES ($1, $2, NOW())`,
        idHistoryState,
        comment)

    if err != nil {
        log.Fatal(err)
    }
    return nil
}

func DBgetComment(id_stage int) ([]models.Stages, error) {
	conn:= db.GetDB()
	if conn == nil {
		return nil, errors.New("DB connection is nil")
	}
   
    rows, err := conn.Query( context.Background(),`
        SELECT c.*
        FROM comments c
        JOIN history_status hs ON c.id_history_state = hs.id_history_status
        WHERE hs.id_stage = $1
    `, id_stage)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var comments []models.Stages
    for rows.Next() {
        var comment models.Stages
        err = rows.Scan(&comment.Id_comment,
            &comment.Id_history_state,
            &comment.Comment,
			&comment.Data_create,)

        if err != nil {
            log.Fatal(err)
        }
        comments = append(comments, comment)
    }
    return comments, nil
}
func DBChengeStatusStage(id_stage int, id_status_stage int, comment string) error {
	conn:= db.GetDB()
	if conn == nil {
		return errors.New("DB connection is nil")
	}
   
    tx, err := conn.Begin( context.Background(),)
    if err != nil {
        return err
    }

    var id_history_status int
    err = tx.QueryRow(context.Background(),`INSERT INTO history_status (id_stage, id_status_stage, data_change_status)
        VALUES ($1, $2, NOW()) RETURNING id_history_status`,
        id_stage,
        id_status_stage).Scan(&id_history_status)

    if err != nil {
        tx.Rollback( context.Background())
        return err
    }

    _, err = tx.Exec( context.Background(),`INSERT INTO comments (id_history_state, comment, date_create_comment)
        VALUES ($1, $2, NOW())`,
        id_history_status,
        comment)

    if err != nil {
        tx.Rollback(context.Background())
        return err
    }

    return tx.Commit( context.Background(),)
}
func DBdeleteFile(id_files int) error {
	conn := db.GetDB()
	

	_, err := conn.Exec( context.Background(),`DELETE FROM files WHERE id_file=$1`, id_files)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func DBdeleteStage(id_stage int) error {
    conn:= db.GetDB()
	if conn == nil {
		return errors.New("DB connection is nil")
	}

    tx, err := conn.Begin( context.Background(),)
    if err != nil {
        return err
    }

    _, err = tx.Exec( context.Background(),`DELETE FROM history_states WHERE id_stage=$1`, id_stage)
    if err != nil {
        tx.Rollback( context.Background(),)
        return err
    }

    _, err = tx.Exec( context.Background(),`DELETE FROM files WHERE id_stage=$1`, id_stage)
    if err != nil {
        tx.Rollback( context.Background(),)
        return err
    }

    _, err = tx.Exec( context.Background(),`DELETE FROM stages WHERE id_stage=$1`, id_stage)
    if err != nil {
        tx.Rollback( context.Background(),)
        return err
    }

    err = tx.Commit( context.Background(),)
    if err != nil {
        return err
    }

    return nil
}

func DBdeleteComment(id_comment int) error {
	conn:= db.GetDB()
	if conn == nil {
		return errors.New("DB connection is nil")
	}
	

	_, err := conn.Exec( context.Background(),`DELETE FROM comments WHERE id_comment=$1`, id_comment)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}
