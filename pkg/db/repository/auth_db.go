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
    err = conn.QueryRow(`SELECT id ,login, password FROM users WHERE login=$1 AND password=$2`, 
                        login, password).Scan(&user.Id_user, &user.Login, &user.Password)

    if err != nil {
        if err == sql.ErrNoRows {
            return 0, errors.New("user not found")
        }
        return 0, err
    }

    return user.Id_user, nil
}

func GetAddmin(id int)(bool,error){
	conn,err:=db.ConnectDB()
	if err!=nil{
		log.Fatal(err)
	}
	var isAdmin bool
	defer conn.Close()
	err=conn.QueryRow(`Select id, admin from users where id=$1`,id).Scan(isAdmin)
	if err !=nil{
		if err==sql.ErrNoRows{
			return false, errors.New("user not found")
		}
		return false, err
	}
	return isAdmin, nil
}