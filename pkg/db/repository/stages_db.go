package db

import (
	"appContract/pkg/db"
	"appContract/pkg/models"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"log"
	"time"

	"github.com/jackc/pgx"
)

func DBgetStageAll() ([]models.Stages, error) {
	conn := db.GetDB()
	if conn == nil {
		return nil, errors.New("DB connection is nil")
	}
	rows, err := conn.Query(context.Background(), `SELECT 
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

    rows, err := conn.Query(context.Background(), `
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
            latest_status.data_change_status,
            contract_user.surname AS contract_surname,
            contract_user.username AS contract_username,
            contract_user.patronymic AS contract_patronymic
        FROM stages s
        JOIN contracts c ON s.id_contract = c.id_contract
        JOIN users u ON s.id_user = u.id_user
        JOIN users contract_user ON c.id_user = contract_user.id_user
        LEFT JOIN LATERAL (
            SELECT 
                hs.id_status_stage,
                hs.data_change_status
            FROM history_status hs
            WHERE hs.id_stage = s.id_stage
            ORDER BY hs.data_change_status DESC
            LIMIT 1
        ) latest_status ON true
        LEFT JOIN status_stages ss ON latest_status.id_status_stage = ss.id_status_stage
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
            &stage.Date_create_start,
            &stage.Date_create_end,
            &stage.Id_contract,
            &stage.Name_contract,
            &stage.Data_contract_create,
            &stage.Id_status_stage,
            &stage.Name_status_stage,
            &stage.Date_change_status,
            &stage.ContractSurname,
            &stage.ContractUsername,
            &stage.ContractPatronymic,
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
            u.patronymic,       
            cu.surname AS contract_surname,     
            cu.username AS contract_username,  
            cu.patronymic AS contract_patronymic,
            hs.data_change_status
        FROM stages s
        JOIN contracts c ON s.id_contract = c.id_contract
        JOIN users cu ON c.id_user = cu.id_user 
        JOIN users u ON s.id_user = u.id_user
        LEFT JOIN LATERAL (
            SELECT 
                id_status_stage,
                data_change_status
            FROM history_status
            WHERE id_stage = s.id_stage
            ORDER BY data_change_status DESC
            LIMIT 1
        ) hs ON true
        LEFT JOIN status_stages ss ON hs.id_status_stage = ss.id_status_stage
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
            &stage.ContractSurname,
            &stage.ContractUsername,
            &stage.ContractPatronymic,
            &stage.Date_change_status,
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
    conn := db.GetDB()
    if conn == nil {
        return models.Stages{}, errors.New("DB connection is nil")
    }

    // Проверяем существование этапа
    var exists bool
    err := conn.QueryRow(context.Background(), 
        "SELECT EXISTS(SELECT 1 FROM stages WHERE id_stage = $1)", 
        stage_id).Scan(&exists)
    if err != nil {
        return models.Stages{}, fmt.Errorf("existence check failed: %v", err)
    }
    if !exists {
        return models.Stages{}, fmt.Errorf("stage with id %d does not exist", stage_id)
    }

    // Основной запрос для получения информации о этапе
    query := `
SELECT 
    s.id_stage,
    s.name_stage,
    s.id_user,
    COALESCE(u.surname, '') as surname,
    COALESCE(u.username, '') as username,
    COALESCE(u.patronymic, '') as patronymic,
    s.description,
    s.date_create_start,
    s.date_create_end,
    s.id_contract,
    COALESCE(cu.surname, '') as contract_surname,
    COALESCE(cu.username, '') as contract_username,
    COALESCE(cu.patronymic, '') as contract_patronymic
FROM stages s
LEFT JOIN users u ON s.id_user = u.id_user
LEFT JOIN contracts c ON s.id_contract = c.id_contract
LEFT JOIN users cu ON c.id_user = cu.id_user
WHERE s.id_stage = $1`

    var stage models.Stages
    err = conn.QueryRow(context.Background(), query, stage_id).Scan(
        &stage.Id_stage,
        &stage.Name_stage,
        &stage.Id_user,
        &stage.Surname,
        &stage.Username,
        &stage.Patronymic,
        &stage.Description,
        &stage.Date_create_start,
        &stage.Date_create_end,
        &stage.Id_contract,
        &stage.ContractSurname,
        &stage.ContractUsername,
        &stage.ContractPatronymic,
    )

    if err != nil {
        return models.Stages{}, fmt.Errorf("query failed: %v", err)
    }

    // Отдельный запрос для получения последнего статуса
    statusQuery := `
SELECT 
    hs.id_status_stage,
    ss.name_status_stage
FROM history_status hs
JOIN status_stages ss ON hs.id_status_stage = ss.id_status_stage
WHERE hs.id_stage = $1
ORDER BY hs.data_change_status DESC
LIMIT 1`

    err = conn.QueryRow(context.Background(), statusQuery, stage_id).Scan(
        &stage.Id_status_stage,
        &stage.Name_status_stage,
    )

    // Если статус не найден, устанавливаем нулевые значения
    if err != nil && err != sql.ErrNoRows {
        return models.Stages{}, fmt.Errorf("status query failed: %v", err)
    }
    if err == sql.ErrNoRows {
        stage.Id_status_stage = 0
        stage.Name_status_stage = ""
    }

    return stage, nil
}
func DBgetFileIDStageID(id_stage int, id_file int) (models.File, error) {
	conn := db.GetDB()
	if conn == nil {
		return models.File{}, errors.New("DB connection is nil")
	}

	var file models.File
	var data []byte
	err := conn.QueryRow(
		context.Background(),
		`SELECT id_file, name_file, data, type_file, id_stage 
         FROM files 
         WHERE id_stage = $1 AND id_file = $2`,
		id_stage,
		id_file,
	).Scan(
		&file.Id_file,
		&file.Name_file,
		&data,
		&file.Type_file,
		&file.Id_stage,
	)

	file.Data = data

	if err != nil {
		if err == pgx.ErrNoRows {
			return models.File{}, fmt.Errorf("file not found")
		}
		return models.File{}, fmt.Errorf("database error: %v", err)
	}

	return file, nil
}

func DBgetFilesStageID(id_stages int) ([]models.File, error) {
	conn := db.GetDB()
	if conn == nil {
		return nil, errors.New("DB connection is nil")
	}

	rows, err := conn.Query(context.Background(), "SELECT * FROM files WHERE Id_stage=$1", id_stages)
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
	conn := db.GetDB()
	if conn == nil {
		return models.StatusStage{}, errors.New("DB connection is nil")
	}

	rows, err := conn.Query(context.Background(), `SELECT * 
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
	conn := db.GetDB()
	if conn == nil {
		return errors.New("DB connection is nil")
	}

	_, err := conn.Exec(context.Background(), `
        INSERT INTO files (
            name_file,
            data,
            type_file,
            id_stage
        ) VALUES ($1, $2, $3, $4)`,
		file.Name_file,
		file.Data,
		file.Type_file,
		file.Id_stage,
	)

	return err
}
func DBaddStage(stage models.Stages) (int, error) {
    conn := db.GetDB()
    if conn == nil {
        return 0, errors.New("DB connection is nil")
    }

    // Начинаем транзакцию
    tx, err := conn.Begin(context.Background())
    if err != nil {
        return 0, fmt.Errorf("failed to begin transaction: %v", err)
    }
    defer tx.Rollback(context.Background())

    // Вставляем запись в stages и получаем ID новой записи
    var stageID int
    err = tx.QueryRow(context.Background(), `
        INSERT INTO stages(
            name_stage,
            id_user,
            description,
            date_create_start,
            date_create_end,
            id_contract
        ) VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id_stage`,
        stage.Name_stage,
        stage.Id_user,
        stage.Description,
        stage.Date_create_start,
        stage.Date_create_end,
        stage.Id_contract).Scan(&stageID)

    if err != nil {
        return 0, fmt.Errorf("failed to insert stage: %v", err)
    }

    // Добавляем начальный статус (1) в историю статусов
    _, err = tx.Exec(context.Background(), `
        INSERT INTO history_status(
            id_stage,
            id_status_stage,
            data_change_status
        ) VALUES ($1, 1, $2)`,
        stageID,
        time.Now().Format("2006-01-02")) // Форматируем дату для поля date

    if err != nil {
        return 0, fmt.Errorf("failed to insert status history: %v", err)
    }

    // Фиксируем транзакцию
    if err := tx.Commit(context.Background()); err != nil {
        return 0, fmt.Errorf("failed to commit transaction: %v", err)
    }

    return stageID, nil
}

func DBaddComment(idStage int, idStatusStage int, comment string, idUser int) error {
	conn := db.GetDB()
	if conn == nil {
		return errors.New("DB connection is nil")
	}

	var idHistoryState int
	err := conn.QueryRow(context.Background(), `SELECT id_history_status FROM history_status
        WHERE id_stage = $1 AND id_status_stage = $2`,
		idStage,
		idStatusStage).Scan(&idHistoryState)

	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(context.Background(), `INSERT INTO comments (id_history_status, comment, id_user, date_create_comment)
        VALUES ($1, $2, $3, NOW())`,
		idHistoryState,
		comment,
		idUser)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func DBgetComment(id_stage int) ([]models.Stages, error) {
	conn := db.GetDB()
	if conn == nil {
		return nil, errors.New("DB connection is nil")
	}

	rows, err := conn.Query(context.Background(), `
        SELECT 
    c.id_comment,
    c.id_history_status,
    c.comment,
    c.date_create_comment,
    c.id_user,
    hs.id_stage,
    hs.id_status_stage,
	ss.name_status_stage,
    u.surname,
    u.username,
    u.patronymic
FROM comments c
JOIN users u ON c.id_user = u.id_user
JOIN history_status hs ON c.id_history_status = hs.id_history_status
JOIN status_stages ss ON hs.id_status_stage = ss.id_status_stage
WHERE hs.id_stage = $1`, id_stage)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Stages
	for rows.Next() {
		var comment models.Stages
		err = rows.Scan(
			&comment.Id_comment,
			&comment.Id_history_status,
			&comment.Comment,
			&comment.Data_create,
			&comment.Id_user,
			&comment.Id_stage,
			&comment.Id_status_stage,
			&comment.Name_status_stage,
			&comment.Surname,
			&comment.Username,
			&comment.Patronymic,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
func DBChengeStatusStage(id_stage int, id_status_stage int, comment string, id_user int) error {
    conn := db.GetDB()
    if conn == nil {
        return errors.New("DB connection is nil")
    }

    tx, err := conn.Begin(context.Background())
    if err != nil {
        return err
    }
    defer func() {
        if err != nil {
            tx.Rollback(context.Background())
        }
    }()

    // 1. Проверяем текущий статус этапа (последняя запись в history_status)
    var currentStatus int
    err = tx.QueryRow(context.Background(), `
        SELECT id_status_stage 
        FROM history_status 
        WHERE id_stage = $1 
        ORDER BY data_change_status DESC 
        LIMIT 1`,
        id_stage).Scan(&currentStatus)

    // Если статус не меняется — можно пропустить (или вернуть ошибку)
    if err == nil && currentStatus == id_status_stage {
        return fmt.Errorf("этап уже имеет статус %d", id_status_stage)
    }

    // 2. Добавляем новую запись в историю статусов
    var id_history_status int
    err = tx.QueryRow(context.Background(), `
        INSERT INTO history_status (id_stage, id_status_stage, data_change_status)
        VALUES ($1, $2, NOW())
        RETURNING id_history_status`,
        id_stage, id_status_stage).Scan(&id_history_status)
    if err != nil {
        return fmt.Errorf("ошибка при добавлении истории статуса: %w", err)
    }

    // 3. Добавляем комментарий (теперь с явным указанием NOW() для date_create_comment)
    _, err = tx.Exec(context.Background(), `
        INSERT INTO comments (id_history_status, comment, date_create_comment, id_user)
        VALUES ($1, $2, NOW(), $3)`,
        id_history_status, comment, id_user)
    if err != nil {
        return fmt.Errorf("ошибка при добавлении комментария: %w", err)
    }

    return tx.Commit(context.Background())
}
func DBdeleteFile(id_files int) error {
	conn := db.GetDB()

	_, err := conn.Exec(context.Background(), `DELETE FROM files WHERE id_file=$1`, id_files)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func DBchangeStage(id_stage int, stage models.Stages) error {
    conn := db.GetDB()
    if conn == nil {
        return errors.New("DB connection is nil")
    }

    _, err := conn.Exec(context.Background(), `
        UPDATE stages SET
            name_stage = $2,
            description = $3,
			id_user = $4,
            date_create_start = $5,
            date_create_end = $6
        WHERE id_stage = $1`,
        id_stage,
        stage.Name_stage,
        stage.Description,
		stage.Id_user,
        stage.Date_create_start,
        stage.Date_create_end,
    )

    if err != nil {
        log.Printf("Failed to update stage: %v", err)
        return err
    }
    return nil
}

func DBdeleteStage(id_stage int) error {
	conn := db.GetDB()
	if conn == nil {
		return errors.New("DB connection is nil")
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	_, err = tx.Exec(context.Background(), `
        DELETE FROM comments 
        WHERE id_history_status IN (
            SELECT id_history_status FROM history_status WHERE id_stage = $1
        )`, id_stage)
	if err != nil {
		return fmt.Errorf("failed to delete comments: %v", err)
	}

	_, err = tx.Exec(context.Background(), `DELETE FROM history_status WHERE id_stage = $1`, id_stage)
	if err != nil {
		return fmt.Errorf("failed to delete history_status: %v", err)
	}

	_, err = tx.Exec(context.Background(), `DELETE FROM files WHERE id_stage = $1`, id_stage)
	if err != nil {
		return fmt.Errorf("failed to delete files: %v", err)
	}

	_, err = tx.Exec(context.Background(), `DELETE FROM stages WHERE id_stage = $1`, id_stage)
	if err != nil {
		return fmt.Errorf("failed to delete stage: %v", err)
	}

	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func DBdeleteComment(id_comment int) error {
	conn := db.GetDB()
	if conn == nil {
		return errors.New("DB connection is nil")
	}

	_, err := conn.Exec(context.Background(), `DELETE FROM comments WHERE id_comment=$1`, id_comment)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}
