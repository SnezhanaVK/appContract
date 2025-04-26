package db

import (
	"appContract/pkg/db"
	"appContract/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

func DBgetUserAll() ([]models.Users, error) {
    conn := db.GetDB()
    if conn == nil {
        return nil, errors.New("connection error")
    }

    rows, err := conn.Query(`SELECT 
        u.id_user, 
        u.surname, 
        u.username, 
        u.patronymic, 
        u.phone, 
        u.photo, 
        u.email, 
        r.id_role, 
        r.name_role 
    FROM 
        users u 
    INNER JOIN 
        user_by_role ubr ON u.id_user = ubr.id_user 
    INNER JOIN 
        roles r ON ubr.id_role = r.id_role
    ORDER BY u.id_user`)
    if err != nil {
        return nil, fmt.Errorf("query error: %v", err)
    }
    defer rows.Close()

    usersMap := make(map[int]*models.Users)
    for rows.Next() {
        var user models.Users
        var role models.Role
        
        err := rows.Scan(
            &user.Id_user,
            &user.Surname,
            &user.Username,
            &user.Patronymic,
            &user.Phone,
            &user.Photo,
            &user.Email,
            &role.Id_role,
            &role.Name_role,
        )
        if err != nil {
            return nil, fmt.Errorf("scan error: %v", err)
        }

        // Если пользователь уже есть в мапе, добавляем только роль
        if existingUser, exists := usersMap[user.Id_user]; exists {
            existingUser.Roles = append(existingUser.Roles, role)
            continue
        }

        // Если пользователя нет в мапе, создаем новую запись
        user.Roles = []models.Role{role}
        usersMap[user.Id_user] = &user
    }

    // Преобразуем map в слайс
    var users []models.Users
    for _, u := range usersMap {
        users = append(users, *u)
    }

    return users, nil
}

func DBgetUserID(user_id int) ([]models.Users, error) {
    conn := db.GetDB()
    if conn == nil {
        return nil, errors.New("connection error")
    }

    rows, err := conn.Query(`
        SELECT 
            u.id_user,
            u.surname,
            u.username,
            u.patronymic,
            u.phone,
            u.photo,
            u.email,
            u.login,
            JSON_AGG(JSON_BUILD_OBJECT('id_role', r.id_role, 'name_role', r.name_role)) AS roles
        FROM users u
        LEFT JOIN user_by_role ubr ON u.id_user = ubr.id_user
        LEFT JOIN roles r ON ubr.id_role = r.id_role
        WHERE u.id_user = $1
        GROUP BY u.id_user
    `, user_id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []models.Users
    for rows.Next() {
        var user models.Users
        var rolesJSON []byte
        
        err = rows.Scan(
            &user.Id_user,
            &user.Surname,
            &user.Username,
            &user.Patronymic,
            &user.Phone,
            &user.Photo,
            &user.Email,
            &user.Login,
            &rolesJSON,
        )
        if err != nil {
            return nil, err
        }
        
        // Декодируем JSON с ролями
        if err := json.Unmarshal(rolesJSON, &user.Roles); err != nil {
            return nil, err
        }
        
        users = append(users, user)
    }
    return users, nil
}

func DBaddUser(user models.Users) error{
	conn:= db.GetDB()
	if conn==nil{
		return errors.New("connection error")
	}

	_,err :=conn.Exec(`
	INSERT INTO users (
	surname, 
	username, 
	patronymic, 
	phone, 
	photo, 
	email, 
	login, 
	password

	)VALUES (
		$1,$2,$3,$4,$5,$6,$7,$8
	)
	`, 
	user.Surname,
	user.Username,
	user.Patronymic,
	user.Phone,
	user.Photo,
	user.Email,
	user.Login,
	user.Password,
)

if err!=nil{
	log.Fatal(err)
}
return nil

}
func DBaddUserAdmin(user models.Users) error{
	conn:= db.GetDB()
	if conn==nil{
		return errors.New("connection error")
	}

	_,err :=conn.Exec(`
        INSERT INTO user_by_role (
            id_user,
            id_role
        ) VALUES (
            $1,
            1
        )
    `, 
        user.Id_user,
    )
    if err != nil {
        log.Fatal(err)
    }

    return nil
}
func DBaddUserMeneger(user models.Users) error{
	conn:= db.GetDB()
	if conn==nil{
		return errors.New("connection error")
	}

	_,err :=conn.Exec(`
        INSERT INTO user_by_role (
            id_user,
            id_role
        ) VALUES (
            $1,
            2
        )
    `, 
        user.Id_user,
    )
    if err != nil {
        log.Fatal(err)
    }

    return nil
}
func DBaddUserUser(user models.Users) error{
	conn:= db.GetDB()
	if conn==nil{
		return errors.New("connection error")
	}

	_,err :=conn.Exec(`
        INSERT INTO user_by_role (
            id_user,
            id_role
        ) VALUES (
            $1,
            3
        )
    `, 
        user.Id_user,
    )
    if err != nil {
        log.Fatal(err)
    }

    return nil
}


func DBchangeUser(user models.Users) error{
	conn:= db.GetDB()
	if conn==nil{
		return errors.New("connection error")
	}

	_,err :=conn.Exec(`
	UPDATE users SET 
	surname=$1,
	username=$2,
	patronymic=$3,
	phone=$4,
	photo=$5,
	email=$6,
	login=$7,
	password=$8
	WHERE id_user=$9
	`, 
	user.Surname,
	user.Username,
	user.Patronymic,
	user.Phone,
	user.Photo,
	user.Email,
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
	conn:= db.GetDB()

	if conn==nil{
		return errors.New("connection error")
	}

	_,err :=conn.Exec(`
	DELETE FROM users WHERE id_user=$1`, user_id)
if err!=nil{
	log.Fatal(err)	
}
return nil
}
