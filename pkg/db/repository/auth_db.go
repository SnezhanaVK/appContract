package db

import (
	db "appContract/pkg/db"
	"appContract/pkg/models"
	"database/sql"
	"errors"
	"log"
)

func Authorize(login string, password string) (int, error) {
    conn, err := db.ConnectDB()
    if err != nil {
        return 0, err
    }
    defer func() {
        if err := conn.Close(); err != nil {
            log.Println(err)
        }
    }()

    var user models.Users
    err = conn.QueryRow(`SELECT id_user ,login, password FROM users WHERE login=$1 AND password=$2`, 
                        login, password).Scan(&user.Id_user, &user.Login, &user.Password)

    if err != nil {
        if err == sql.ErrNoRows {
            return 0, errors.New("user not found")
        }
        return 0, err
    }

    return user.Id_user, nil
}

func GetAddmin(id int) (bool, error) {
    conn, err := db.ConnectDB()
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    var isAdmin sql.NullBool
    err = conn.QueryRow(`SELECT admin FROM users WHERE id_user=$1`, id).Scan(&isAdmin)
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
	conn,err:=db.ConnectDB()
	if err!=nil{
		return err
	}
	defer conn.Close()
	_,err=conn.Exec(`UPDATE users SET password=$1 WHERE login=$2`,password,login)
	if err!=nil{
		return err
	}
	return nil
}

func GetUser(login string) (models.Users, error) {
    conn, err := db.ConnectDB()
    if err != nil {
        log.Println(err)
        return models.Users{}, err
    }
    defer conn.Close()

    var user models.Users

    err = conn.QueryRow(`SELECT id_user, email, login, password FROM users WHERE login = $1`, login).Scan(
        &user.Id_user,
        &user.Email,
        &user.Login,
        &user.Password,
    )

    if err != nil {
        log.Println(err)
        if err == sql.ErrNoRows {
            return models.Users{}, errors.New("Пользователь не найден")
        }
        return models.Users{}, err
    }

    log.Println("User found:", user)
    return user, nil
}