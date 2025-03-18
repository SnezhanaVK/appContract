package handlers

import (
	db "appContract/pkg/db/repository"
	"appContract/pkg/models"
	"encoding/json"
	"net/http"
	"strconv"
)

// Users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
    if r.Method !=http.MethodGet{
        http.Error(w,"Invalid request method GetAllUsers",http.StatusBadRequest)
        return
    }
    users, err := db.DBgetUserAll()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(users)
}

func GetUserID(w http.ResponseWriter, r *http.Request) {
    if r.Method !=http.MethodGet{
        http.Error(w,"Invalid request method GetUserID",http.StatusBadRequest)
        return
    }
    userId, err:= strconv.Atoi(r.URL.Query().Get("user_id"))
    if err != nil {
        http.Error(w, "Invalid user_id", http.StatusBadRequest) 
        return
    }

    user, err := db.DBgetUserID(userId)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(user)
    
}

func PostCreateUser(w http.ResponseWriter, r *http.Request) {
    var user models.Users
    err:=json.NewDecoder(r.Body).Decode(&user)
    if err!=nil{
        http.Error(w,"Invalid request body PostCreateUser",http.StatusBadRequest)
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
    if r.Method !=http.MethodPut{
        http.Error(w,"Invalid request method PutUpdateUser",http.StatusBadRequest)
        return
    }
    var user models.Users
    err:=json.NewDecoder(r.Body).Decode(&user)
    if err!=nil{
        http.Error(w,"Invalid request body",http.StatusBadRequest)
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
    if r.Method !=http.MethodDelete{
    http.Error(w,"Invalid request method DeleteUser",http.StatusBadRequest)
    return
    }

    userId, err:= strconv.Atoi(r.URL.Query().Get("user_id"))
    if err !=nil {
        http.Error(w, "Invalid user_id", http.StatusBadRequest)
        return
    }
    err =db.DBdeleteUser(userId)
    if err !=nil{
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})


}
