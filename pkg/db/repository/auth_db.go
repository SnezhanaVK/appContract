package db

import (
	db "appContract/pkg/db"
	"appContract/pkg/models"
	"context"
	"database/sql"

	//"encoding/json"
	"errors"
	//"fmt"
	"log"

	"github.com/jackc/pgx"
	"github.com/lib/pq"
)

// Изменения в SQL-запросе
func Authorize(login string, password string) (*models.Users, error) {
    conn := db.GetDB()
    if conn == nil {
        return nil, errors.New("connection error")
    }

    query := `
        SELECT 
            u.id_user,
            u.email,
            u.login,
            u.password,
            ARRAY_AGG(r.id_role) AS roles  
        FROM users u
        LEFT JOIN user_by_role ubr ON u.id_user = ubr.id_user
        LEFT JOIN roles r ON ubr.id_role = r.id_role
        WHERE u.login = $1 AND u.password = $2
        GROUP BY u.id_user
    `

    var user models.Users
    var roles []int

    // Используем pq.Int32Array для обработки NULL
    var pgRoles pq.Int32Array
    
    err := conn.QueryRow(context.Background(),query, login, password).Scan(
        &user.Id_user,
        &user.Email,
        &user.Login,
        &user.Password,
        &pgRoles,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, errors.New("user not found")
        }
        return nil, err
    }

    // Преобразование pq.Int32Array в []int
    roles = make([]int, len(pgRoles))
    for i, v := range pgRoles {
        roles[i] = int(v)
    }

    // Обработка NULL массива
    if len(roles) == 1 && roles[0] == 0 {
        roles = []int{}
    }

    user.Admin = false
    user.Manager = false
    for _, role := range roles {
        if role == 1 {
            user.Admin = true
        }
        if role == 2 {
            user.Manager = true
        }
    }

    return &user, nil
}

func GetAddmin(id int) (bool, error) {
	conn := db.GetDB()
	if conn == nil {
		return false, errors.New("connection error")
	}

	var isAdmin sql.NullBool
	err := conn.QueryRow( context.Background(),`SELECT admin FROM users WHERE id_user=$1`, id).Scan(&isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("user not found")
		}
		return false, err
	}

	if !isAdmin.Valid {
		return false, errors.New("admin value is null")
	}

	return isAdmin.Bool, nil
}

func ChangePassword(email string, password string) error {
	if password == "" {
		return errors.New("password is required")
	}
	conn := db.GetDB()
	if conn == nil {
		return errors.New("connection error")
	}

	result, err := conn.Exec( context.Background(),
		`UPDATE users SET password = $1 WHERE email = $2`,
		password,
		email,
	)
	if err != nil {
		return err
	}

	// Проверяем, была ли обновлена хотя бы одна строка
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func GetUser(email string) (models.Users, error) {
	conn := db.GetDB()
	if conn == nil {
		return models.Users{}, errors.New("connection error")
	}

	var user models.Users

	err := conn.QueryRow( context.Background(),`SELECT id_user, email, login, password FROM users WHERE email = $1`, email).Scan(
		&user.Id_user,
		&user.Email,
		&user.Login,
		&user.Password,
	)

	if err != nil {
		log.Println(err)
		if err == pgx.ErrNoRows {
			return models.Users{}, errors.New("User not found")
		}
		return models.Users{}, err
	}

	log.Println("User found:", user)
	return user, nil
}

func GetUserByEmail(email string) (models.Users, error) {
	conn := db.GetDB()
	if conn == nil {
		return models.Users{}, errors.New("connection error")
	}

	var user models.Users

	err := conn.QueryRow( context.Background(),`SELECT id_user, email, login, password FROM users WHERE email = $1`, email).Scan(
		&user.Id_user,
		&user.Email,
		&user.Login,
		&user.Password,
	)

	if err != nil {
		log.Println(err)
		if err == pgx.ErrNoRows {
			return models.Users{}, errors.New("Пользователь не найден")
		}
		return models.Users{}, err
	}

	log.Println("User found:", user)
	return user, nil
}
