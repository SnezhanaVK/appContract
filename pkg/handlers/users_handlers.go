package handlers

import (
	db "appContract/pkg/db/repository"
	"appContract/pkg/models"
	servis "appContract/pkg/services"

	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    dbUsers, err := db.DBgetUserAll()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Создаем кастомную структуру ответа
    type RoleResponse struct {
        IdRole   int    `json:"id_role"`
        NameRole string `json:"name_role"`
    }

    type UserResponse struct {
        IdUser     int            `json:"id_user"`
        Surname    string         `json:"surname"`
        Username   string         `json:"username"`
        Patronymic string         `json:"patronymic,omitempty"`
        Phone      string         `json:"phone"`
        Photo      string         `json:"photo,omitempty"`
        Email      string         `json:"email"`
        Roles      []RoleResponse `json:"roles"`
    }

    // Преобразуем исходные данные
    response := make([]UserResponse, 0, len(dbUsers))
    for _, user := range dbUsers {
        roles := make([]RoleResponse, 0, len(user.Roles))
        for _, role := range user.Roles {
            roles = append(roles, RoleResponse{
                IdRole:   role.Id_role,
                NameRole: role.Name_role,
            })
        }

        response = append(response, UserResponse{
            IdUser:     user.Id_user,
            Surname:    user.Surname,
            Username:   user.Username,
            Patronymic: user.Patronymic,
            Phone:      user.Phone,
            Photo:      user.Photo,
            Email:      user.Email,
            Roles:      roles,
        })
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, "Error creating response", http.StatusInternalServerError)
    }
}




func GetUserID(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusBadRequest)
        return
    }

    vars := mux.Vars(r)
    userId := vars["userID"]
    if userId == "" {
        http.Error(w, "Invalid user_id", http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(userId)
    if err != nil {
        http.Error(w, "Invalid user_id", http.StatusBadRequest)
        return
    }

    users, err := db.DBgetUserID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var response []map[string]interface{}
    for _, user := range users {
        userResponse := map[string]interface{}{
            "id_user":    user.Id_user,
            "surname":    user.Surname,
            "username":   user.Username,
            "patronymic": user.Patronymic,
            "phone":      user.Phone,
            "photo":      user.Photo,
            "email":      user.Email,
            "roles":      user.Roles,
        }
        response = append(response, userResponse)
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func PostCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}
	var user models.Users
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request body PostCreateUser", http.StatusBadRequest)
        return
    }

    err = db.DBaddUser(user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

func PostAddUserRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method PostAddUserRole", http.StatusBadRequest)
		return
	}
	var user models.Users

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body PostAddUserRole", http.StatusBadRequest)
		return
	}
	err = servis.AddRole(user.Id_user, user.Id_role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User role added successfully"})
}


func PutUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method PutUpdateUser", http.StatusBadRequest)
		return
	}
	var user models.Users
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = db.DBchangeUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method DeleteUser", http.StatusBadRequest)
		return
	}


	vars:=mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	err = db.DBdeleteUser(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})

}
