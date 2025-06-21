package db

import (
	"appContract/pkg/db"
	"appContract/pkg/models"
	"appContract/pkg/utils"
	"context"
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

	rows, err := conn.Query(context.Background(), `SELECT 
        u.id_user, 
        u.surname, 
        u.username, 
        u.patronymic, 
        u.phone, 
        u.email,
		u.login, 
        r.id_role, 
        r.name_role 
    FROM 
        users u 
    LEFT JOIN 
        user_by_role ubr ON u.id_user = ubr.id_user 
    LEFT JOIN 
        roles r ON ubr.id_role = r.id_role
    ORDER BY u.id_user`)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	usersMap := make(map[int]*models.Users)
	for rows.Next() {
		var user models.Users
		var roleID *int
		var roleName *string

		err := rows.Scan(
			&user.Id_user,
			&user.Surname,
			&user.Username,
			&user.Patronymic,
			&user.Phone,
			&user.Email,
			&user.Login,
			&roleID,
			&roleName,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		if existingUser, exists := usersMap[user.Id_user]; exists {
			if roleID != nil {
				role := models.Role{
					Id_role:   *roleID,
					Name_role: *roleName,
				}
				existingUser.Roles = append(existingUser.Roles, role)
			}
			continue
		}

		user.Roles = []models.Role{}
		if roleID != nil {
			role := models.Role{
				Id_role:   *roleID,
				Name_role: *roleName,
			}
			user.Roles = append(user.Roles, role)
		}
		usersMap[user.Id_user] = &user
	}

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

	rows, err := conn.Query(context.Background(), `
        SELECT 
            u.id_user,
            u.surname,
            u.username,
            u.patronymic,
            u.phone,
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
			&user.Email,
			&user.Login,
			&rolesJSON,
		)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(rolesJSON, &user.Roles); err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

func DBaddUser(user models.Users, password string)  error {
	conn := db.GetDB()
	if conn == nil {
		return errors.New("connection error")
	}

	salt, err := utils.GenerateSalt(16)
	if err != nil {
		return fmt.Errorf("failed to generate salt: %v", err)
	}

	hashedPassword, err := utils.HashPassword(password, salt)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	_, err = conn.Exec(context.Background(), `
    INSERT INTO users (
        surname, 
        username, 
        patronymic, 
        phone, 
        email, 
        login, 
        password_hash,
        salt,
        password_algorithm
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9
    )`,
		user.Surname,
		user.Username,
		user.Patronymic,
		user.Phone,
		user.Email,
		user.Login,
		hashedPassword,
		salt,
		"bcrypt",
	)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}
	return  nil
}
func DBgetUserId(login string)(int, error) {
	conn := db.GetDB()
	if conn == nil {
		return 0, errors.New("connection error")
	}
	var id int
	err := conn.QueryRow(context.Background(), "SELECT id_user FROM users WHERE login= $1", login).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func DBAddUserRole(user models.Users, roleID int) error {
	conn := db.GetDB()
	if conn == nil {
		return errors.New("connection error")
	}
	var roleName string
	switch roleID {
	case 1:
		roleName = "admin"
	case 2:
		roleName = "manager"
	default:
		return fmt.Errorf("unknown role ID: %d", roleID)
	}

	var exists bool
	err := conn.QueryRow(context.Background(), `
        SELECT EXISTS(
            SELECT 1 FROM user_by_role 
            WHERE id_user = $1 AND id_role = $2
        )
    `, user.Id_user, roleID).Scan(&exists)

	if err != nil {
		log.Printf("Error checking existing %s role: %v", roleName, err)
		return fmt.Errorf("failed to check %s role: %v", roleName, err)
	}

	if exists {
		return fmt.Errorf("user already has %s role (id_role=%d)", roleName, roleID)
	}

	_, err = conn.Exec(context.Background(), `
        INSERT INTO user_by_role (
            id_user,
            id_role
        ) VALUES (
            $1,
            $2
        )
    `, user.Id_user, roleID)

	if err != nil {
		log.Printf("Error adding %s role: %v", roleName, err)
		return fmt.Errorf("failed to add %s role: %v", roleName, err)
	}

	return nil
}

func DBaddUserAdmin(user models.Users) error {
	return DBAddUserRole(user, 1)
}

func DBaddUserMeneger(user models.Users) error {
	return DBAddUserRole(user, 2)
}

func DBRemoveUserRole(user models.Users, roleID int) error {
	conn := db.GetDB()
	if conn == nil {
		return errors.New("connection error")
	}

	// Определяем название роли для сообщений об ошибках
	var roleName string
	switch roleID {
	case 1:
		roleName = "admin"
	case 2:
		roleName = "manager"
	default:
		return fmt.Errorf("unknown role ID: %d", roleID)
	}

	// Проверяем, есть ли указанная роль у пользователя
	var exists bool
	err := conn.QueryRow(context.Background(), `
        SELECT EXISTS(
            SELECT 1 FROM user_by_role 
            WHERE id_user = $1 AND id_role = $2
        )`,
		user.Id_user,
		roleID,
	).Scan(&exists)

	if err != nil {
		log.Printf("Error checking %s role existence: %v", roleName, err)
		return fmt.Errorf("failed to check %s role: %v", roleName, err)
	}

	if !exists {
		return fmt.Errorf("user doesn't have %s role (id_role=%d)", roleName, roleID)
	}

	tag, err := conn.Exec(context.Background(), `
        DELETE FROM user_by_role
        WHERE id_user = $1 AND id_role = $2`,
		user.Id_user,
		roleID,
	)

	if err != nil {
		log.Printf("Error removing %s role: %v", roleName, err)
		return fmt.Errorf("failed to remove %s role: %v", roleName, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no %s role was removed (id_user=%d)", roleName, user.Id_user)
	}

	return nil
}
func DBRemoveUserAdmin(user models.Users) error {
	return DBRemoveUserRole(user, 1)
}

func DBRemoveUserMeneger(user models.Users) error {
	return DBRemoveUserRole(user, 2)
}

func DBgetUserRoles(user_id int) ([]models.Role, error) {
	conn := db.GetDB()
	if conn == nil {
		return nil, errors.New("connection error")
	}

	rows, err := conn.Query(context.Background(), `
        SELECT 
            r.id_role, 
            r.name_role 
        FROM 
            user_by_role ubr 
        INNER JOIN 
            roles r ON ubr.id_role = r.id_role 
        WHERE 
            ubr.id_user = $1`, user_id)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.Id_role, &role.Name_role)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func DBchangeUser(user models.Users) error {
	conn := db.GetDB()
	if conn == nil {
		return errors.New("database connection error")
	}

	if user.Id_user == 0 {
		return errors.New("user ID is required")
	}

	result, err := conn.Exec(context.Background(), `
    UPDATE users SET 
        surname = COALESCE(NULLIF($1, ''), surname),
        username = COALESCE(NULLIF($2, ''), username),
        patronymic = COALESCE(NULLIF($3, ''), patronymic),
        phone = COALESCE(NULLIF($4, ''), phone),
		login = COALESCE(NULLIF($5, ''), login),
		email = COALESCE(NULLIF($6, ''), email)
   	    WHERE id_user = $7
    `,
		user.Surname,
		user.Username,
		user.Patronymic,
		user.Phone,	
		user.Login,
		user.Email,
		user.Id_user,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no rows were updated - user not found or no changes made")
	}

	return nil
}
func DBdeleteUser(user_id int) error {
	conn := db.GetDB()
	if conn == nil {
		return errors.New("connection error")
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("transaction begin error: %v", err)
	}
	defer tx.Rollback(context.Background())

	var contractCount int
	err = tx.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM contracts WHERE id_user = $1", user_id).Scan(&contractCount)
	if err != nil {
		return fmt.Errorf("contracts check error: %v", err)
	}
	if contractCount > 0 {
		return fmt.Errorf("cannot delete user - user has %d associated contracts", contractCount)
	}

	var stageCount int
	err = tx.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM stages WHERE id_user = $1", user_id).Scan(&stageCount)
	if err != nil {
		return fmt.Errorf("stages check error: %v", err)
	}
	if stageCount > 0 {
		return fmt.Errorf("cannot delete user - user has %d associated stages", stageCount)
	}

	_, err = tx.Exec(context.Background(),
		"DELETE FROM user_photos WHERE id_user = $1", user_id)
	if err != nil {
		return fmt.Errorf("user_photos delete error: %v", err)
	}

	_, err = tx.Exec(context.Background(),
		"DELETE FROM notification_settings_by_user WHERE id_user = $1", user_id)
	if err != nil {
		return fmt.Errorf("notification_settings_by_user delete error: %v", err)
	}

	_, err = tx.Exec(context.Background(),
		"DELETE FROM user_by_role WHERE id_user = $1", user_id)
	if err != nil {
		return fmt.Errorf("user_by_role delete error: %v", err)
	}

	result, err := tx.Exec(context.Background(),
		"DELETE FROM users WHERE id_user = $1", user_id)
	if err != nil {
		return fmt.Errorf("users delete error: %v", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", user_id)
	}

	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("transaction commit error: %v", err)
	}

	return nil
}
