package handlers

import (
	db "appContract/pkg/db/repository"
	"encoding/json"
	"net/http"
)

type LoginStruct struct{
	Login    string `json:"login"`
	Password string `json:"password"`
}
//Login
func Login(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusBadRequest)
        return
    }

    var authRequest LoginStruct
    err := json.NewDecoder(r.Body).Decode(&authRequest)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    login := authRequest.Login
    password := authRequest.Password

    if login == "" || password == "" {
        http.Error(w, "Login and password are required", http.StatusBadRequest)
        return
    }

    if db.Authorize(login, password) {
        w.Write([]byte("Authorized"))
    } else {
        http.Error(w, "Invalid login or password", http.StatusUnauthorized)
    }
}

func ForgotPassword(w http.ResponseWriter, r *http.Request){
}










