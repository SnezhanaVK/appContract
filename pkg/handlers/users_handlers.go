package handlers

import (
	db "appContract/pkg/db/repository"
	"appContract/pkg/models"
	"appContract/pkg/service"
	"appContract/pkg/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

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
		Email      string         `json:"email"`
		Login      string         `json:"login"`
		Roles      []RoleResponse `json:"roles"`
	}

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
			Email:      user.Email,
			Login:      user.Login,
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
			"login":      user.Login,

			"email": user.Email,
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
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	id_user,err:=db.DBgetUserId(user.Login)
	log.Println(id_user,err)
	if id_user!=0{
		http.Error(w, "User with this login already exists", http.StatusBadRequest)
		return
	}
	
	emailSender := utils.NewDefaultEmailSender()
	if err := service.CreateUser(user, emailSender); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully. Account credentials sent to email.",
		"ID":       strconv.Itoa(id_user),
	})
}

func PostAddRoleAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method PostAddRoleAdmin", http.StatusBadRequest)
		return
	}
	var user models.Users
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
	user.Id_user = id

	err = db.DBaddUserAdmin(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User role added successfully"})
}

func PostAddRoleManager(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method PostAddRoleMeneger", http.StatusBadRequest)
		return
	}
	var user models.Users
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
	user.Id_user = id

	err = db.DBaddUserMeneger(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User role added successfully"})
}

func DeleteRemoveRoleAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method DeleteRemoveRoleAdmin", http.StatusBadRequest)
		return
	}
	var user models.Users
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
	user.Id_user = id

	err = db.DBRemoveUserRole(user, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User role removed successfully"})
}

func DeleteRemoveRoleManager(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method DeleteRemoveRoleMeneger", http.StatusBadRequest)
		return
	}
	var user models.Users
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
	user.Id_user = id

	err = db.DBRemoveUserRole(user, 2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User role removed successfully"})
}

func GetUserRoles(w http.ResponseWriter, r *http.Request) {
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

	roles, err := db.DBgetUserRoles(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(roles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func PutUpdateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method, PUT expected", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	if userIDStr == "" {
		http.Error(w, "User ID is required in URL", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	var user models.Users
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	user.Id_user = userID

	if err := db.DBchangeUser(user); err != nil {
		if strings.Contains(err.Error(), "no rows were updated") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User updated successfully",
		"user_id": userID,
	})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	err = db.DBdeleteUser(userId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("User with id %d deleted successfully", userId),
	})
}
