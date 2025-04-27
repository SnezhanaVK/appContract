package service

// //db "appContract/pkg/db/repository"
// import (
// 	"appContract/pkg/db"
// 	"appContract/pkg/models"
// 	"appContract/pkg/utils"
// 	"errors"
// 	"log"
// 	"time"
// )
// func SendingCode(models.Users) (string, error) {

// 	conn := db.GetDB()
// 	if conn == nil {
// 		return "", errors.New("DB connection is nil")
// 	}

// 	var code string
// 	err := conn.QueryRow(`SELECT code FROM users WHERE email = $1`, email).Scan(&c)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return "", errors.New("user not found")
// 		}
// 		return "", err
// 	}
// }
// code




