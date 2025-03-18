package db

import (
	"appContract/pkg/db"
	"appContract/pkg/models"
	"log"
)

func DBgetUserAll() ([]models.Users, error) {
	// соединение с бд
	conn, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// запрос к бд
	rows, err := conn.Query(`SELECT 
							id_user, 
							surname, 
							username, 
							patronymic, 
							phone, 
							photo, 
							email 
							FROM users`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// обработка результата
	var users []models.Users
	for rows.Next() {
		var user models.Users
		err = rows.Scan(&user.Id_user, 
						&user.Surname, 
						&user.Username, 
						&user.Patronymic, 
						&user.Phone, 
						&user.Photo, 
						&user.Email, 
		) 
						
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	return users, nil
}

func DBgetUserID( user_id int)([]models.Users, error){
conn, err:= db.ConnectDB()
if err!=nil{
	log.Fatal(err)
}
defer conn.Close()

rows, err := conn.Query(`SELECT 
  u.id_user, 
  u.surname, 
  u.username, 
  u.patronymic, 
  u.phone, 
  u.photo, 
  u.email,
  u.login,
  u.notification_id,
  n.variant_notification
FROM users u
JOIN notifications n ON u.notification_id = n.id_notification
							WHERE id_user=$1`,user_id)
if err != nil {
	log.Fatal(err)
}
defer rows.Close()

// обработка результата
var users []models.Users
for rows.Next() {
	var user models.Users
	err = rows.Scan(&user.Id_user, 
					&user.Surname, 
					&user.Username, 
					&user.Patronymic, 
					&user.Phone, 
					&user.Photo,
					&user.Email,
					&user.Login,
					&user.Notification_id,
					&user.Variant_notification,)
	if err != nil {
		log.Fatal(err)
	}
	users = append(users, user)
}
return users, nil
}

func DBaddUser(user models.Users) error{
	conn, err:= db.ConnectDB()
	if err!=nil{
		log.Fatal(err)
	}
	defer conn.Close()

	_,err =conn.Exec(`
	INSERT INTO users (
	id_user, 
	surname, 
	username, 
	patronymic, 
	phone, 
	photo, 
	email, 
	role_id, 
	notification_id, 
	admin, 
	login, 
	password
	)VALUES (
		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11
	)
	`, 
	user.Surname,
	user.Username,
	user.Patronymic,
	user.Phone,
	user.Photo,
	user.Email,
	user.Role_id,
	user.Notification_id,
	user.Admin,
	user.Login,
	user.Password,
)
if err!=nil{
	log.Fatal(err)
}
return nil
}

func DBchangeUser(user models.Users) error{
	conn, err:= db.ConnectDB()
	if err!=nil{
		log.Fatal(err)
	}
	defer conn.Close()

	_,err =conn.Exec(`
	UPDATE users SET 
	surname=$1,
	username=$2,
	patronymic=$3,
	phone=$4,
	photo=$5,
	email=$6,
	role_id=$7,
	notification_id=$8,
	admin=$9,
	login=$10,
	password=$11
	WHERE id_user=$12
	`, 
	user.Surname,
	user.Username,
	user.Patronymic,
	user.Phone,
	user.Photo,
	user.Email,
	user.Role_id,
	user.Notification_id,
	user.Admin,
	user.Login,
	user.Password,
	user.Id_user,
)
if err!=nil{
	log.Fatal(err)
}
return nil
}

func DBdeleteUser(user_id int) error{
	conn, err:= db.ConnectDB()
	if err!=nil{
		log.Fatal(err)
	}
	defer conn.Close()

	_,err =conn.Exec(`
	DELETE FROM users WHERE id_user=$1`, user_id)
if err!=nil{
	log.Fatal(err)	
}
return nil
}
