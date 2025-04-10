package handlers

import (
	db "appContract/pkg/db/repository"
	"appContract/pkg/models"

	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method GetAllUsers", http.StatusBadRequest)
		return
	}
	users, err := db.DBgetUserAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var userResponses []map[string]interface{}
	for _, user := range users {
		userResponse := map[string]interface{}{
			"id_user":    user.Id_user,
			"surname":    user.Surname,
			"username":   user.Username,
			"patronymic": user.Patronymic,
			"phone":      user.Phone,
			"photo":      user.Photo,
			"email":      user.Email,
		}
		userResponses = append(userResponses, userResponse)
	}
	data, err := json.Marshal(userResponses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func GetUserID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method GetUserID", http.StatusBadRequest)
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

	user, err := db.DBgetUserID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userResponses []map[string]interface{}
	for _, user := range user {
		userResponse := map[string]interface{}{
			"id_user":    user.Id_user,
			"surname":    user.Surname,
			"username":   user.Username,
			"patronymic": user.Patronymic,
			"phone":      user.Phone,
			"photo":      user.Photo,
			"email":      user.Email,
		}
		userResponses = append(userResponses, userResponse)
	}
	data, err := json.Marshal(userResponses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
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
