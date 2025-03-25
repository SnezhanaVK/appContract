package db

import (
	db "appContract/pkg/db"
	"appContract/pkg/models"
	"database/sql"
	"log"
)

func Authorize(login string, password string) bool {
		conn,err:=db.ConnectDB()
		if err!=nil{
			log.Fatal(err)
		}
		defer conn.Close()
		var user models.Users
		err = conn.QueryRow(`SELECT login, password FROM users WHERE login=$1 AND password=$2`, 
							login, password).Scan(&user.Login, &user.Password)

		if err !=nil{
			if err== sql.ErrNoRows{
				return false
			}
			log.Fatal(err)
		}
		 

		return true
	}