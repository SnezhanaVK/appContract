package handlers

import (
	db "appContract/pkg/db/repository"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type LoginStruct struct {
	ID       int    `json:"id"`
    Admin    bool   `json:admin`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

var jwtKey = []byte("secretkey")

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

    _, err = db.Authorize(login, password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    // Создаем токен авторизации
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "login": login,
        "exp":   time.Now().Add(time.Hour * 72).Unix(),
    })

    // Подписываем токен
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Возвращаем токен
    json.NewEncoder(w).Encode(Token{Token: tokenString})
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
    tokenString := r.Header.Get("Authorization")
    if tokenString == "" {
        http.Error(w, "Token is required", http.StatusBadRequest)
        return
    }

    token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }

    login := claims["login"].(string)
    if login == "" {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }

    id, err := db.Authorize(login, "")
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    isAdmin, err := db.GetAddmin(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if isAdmin {
        w.Write([]byte("Authorized as admin"))
    } else {
        w.Write([]byte("Authorized as user"))
    }
}

func PutForgotPassword(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPut {
        http.Error(w, "Invalid request method", http.StatusBadRequest)
        return 
    }
    var authRequest LoginStruct
    err:=json.NewDecoder(r.Body).Decode(&authRequest)
    if err!=nil{
        http.Error(w,"Invalid request body PutChangePassword",http.StatusBadRequest)
        return
    }
    if authRequest.ID==0 || authRequest.Password==""{
        http.Error(w,"Invalid request body PutChangePassword",http.StatusBadRequest)
        return
    }
    err=db.ChangePassword(authRequest.ID,authRequest.Password)
    if err!=nil{
        http.Error(w,err.Error(),http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})

    

}   