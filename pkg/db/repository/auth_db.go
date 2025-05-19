package db

// auth_db.go

import (
	db "appContract/pkg/db"
	"appContract/pkg/models"
	"appContract/pkg/utils"
	"context"
	"database/sql"
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

	userQuery := `
        SELECT 
            id_user,
            surname,
            username,
            patronymic,
            phone,
            email,
            login,
            password_hash,
            salt
        FROM users 
        WHERE login = $1
    `

	var user models.Users
	var passwordHash string
	var salt string

	err := conn.QueryRow(context.Background(), userQuery, login).Scan(
		&user.Id_user,
		&user.Surname,
		&user.Username,
		&user.Patronymic,
		&user.Phone,
		&user.Email,
		&user.Login,
		&passwordHash,
		&salt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if !utils.VerifyPassword(passwordHash, password, salt) {
		return nil, errors.New("invalid password")
	}

	rolesQuery := `
        SELECT r.id_role, r.name_role 
        FROM user_by_role ubr
        JOIN roles r ON ubr.id_role = r.id_role
        WHERE ubr.id_user = $1
    `

	rows, err := conn.Query(context.Background(), rolesQuery, user.Id_user)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %v", err)
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.Id_role, &role.Name_role); err != nil {
			return nil, fmt.Errorf("failed to scan role: %v", err)
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("roles query error: %v", err)
	}

	user.Roles = roles

	if len(user.Roles) > 0 {
		user.Id_role = user.Roles[0].Id_role
		user.Name_role = user.Roles[0].Name_role
	} else {
		user.Id_role = 0
		user.Name_role = ""
	}

	for _, role := range user.Roles {
		switch role.Id_role {
		case 1:
			user.Admin = true
		case 2:
			user.Manager = true
		}
	}

	return &user, nil
}

func GetAdmin(id int) (bool, error) {
	conn := db.GetDB()
	if conn == nil {
		return false, errors.New("connection error")
	}

	var isAdmin sql.NullBool
	err := conn.QueryRow(context.Background(), `SELECT admin FROM users WHERE id_user=$1`, id).Scan(&isAdmin)
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

func ChangePassword(email string, newPassword string) error {
	if newPassword == "" {
		return errors.New("password is required")
	}

	conn := db.GetDB()
	if conn == nil {
		return errors.New("connection error")
	}

	// Генерируем новую соль и хеш
	salt, err := utils.GenerateSalt(16)
	if err != nil {
		return fmt.Errorf("failed to generate salt: %v", err)
	}

	hashedPassword, err := utils.HashPassword(newPassword, salt)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	result, err := conn.Exec(context.Background(),
		`UPDATE users 
         SET password_hash = $1, 
             salt = $2,
             password_updated_at = CURRENT_TIMESTAMP
         WHERE email = $3`,
		hashedPassword,
		salt,
		email,
	)

	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
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

	err := conn.QueryRow(context.Background(), `SELECT id_user, email, login FROM users WHERE email = $1`, email).Scan(
		&user.Id_user,
		&user.Email,
		&user.Login,
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

	err := conn.QueryRow(context.Background(), `SELECT id_user, email, login  FROM users WHERE email = $1`, email).Scan(
		&user.Id_user,
		&user.Email,
		&user.Login,
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
