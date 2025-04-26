package db

import (
	db "appContract/pkg/db"
	"appContract/pkg/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx"
)

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
            COALESCE(json_agg(json_build_object(
                'id_role', r.id_role, 
                'name_role', r.name_role
            )) FILTER (WHERE r.id_role IS NOT NULL), '[]'::json) AS roles
        FROM users u
        LEFT JOIN user_by_role ubr ON u.id_user = ubr.id_user
        LEFT JOIN roles r ON ubr.id_role = r.id_role
        WHERE u.login = $1 AND u.password = $2
        GROUP BY u.id_user
    `

    var user models.Users
    var rolesJSON []byte

    err := conn.QueryRow(query, login, password).Scan(
        &user.Id_user,
        &user.Email,
        &user.Login,
        &user.Password,
        &rolesJSON,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, errors.New("user not found")
        }
        return nil, err
    }

    // Декодируем роли
    if err := json.Unmarshal(rolesJSON, &user.Roles); err != nil {
        return nil, fmt.Errorf("failed to decode roles: %v", err)
    }

    // Заполняем первую роль для совместимости
    if len(user.Roles) > 0 {
        user.Id_role = user.Roles[0].Id_role
        user.Name_role = user.Roles[0].Name_role
    } else {
        user.Id_role = 0
        user.Name_role = ""
    }

    return &user, nil
}

func GetAddmin(id int) (bool, error) {
    conn:= db.GetDB()
    if conn==nil{
        return false, errors.New("connection error")
    }

    var isAdmin sql.NullBool
    err := conn.QueryRow(`SELECT admin FROM users WHERE id_user=$1`, id).Scan(&isAdmin)
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


func ChangePassword(login string, password string)error{
	if password==""{
		return errors.New("password is required")
	}
	conn:=db.GetDB()
    if conn==nil{
        return errors.New("connection error")
    }


	_,err:=conn.Exec(`UPDATE users SET password=$1 WHERE login=$2`,password,login)
	if err!=nil{
		return err
	}
	return nil
}

func GetUser(login string) (models.Users, error) {
    conn:= db.GetDB() 
    if conn==nil{
        return models.Users{}, errors.New("connection error")
    }

    var user models.Users


    err:= conn.QueryRow(`SELECT id_user, email, login, password FROM users WHERE login = $1`, login).Scan(
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