package handlers

import (
	db "appContract/pkg/db/repository"
	"appContract/pkg/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)
type Token struct {
    Token     string        `json:"token"`
    Id_user   int           `json:"id_user"`
    Roles     []models.Role `json:"roles"`              
}

var jwtKey = []byte("secretkey")

func Login(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusBadRequest)
        return
    }

    var user *models.Users
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    login := user.Login
    password := user.Password
    if login == "" || password == "" {
        http.Error(w, "Login and password are required", http.StatusBadRequest)
        return
    }
    
    authUser, err := db.Authorize(login, password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    // Создаем JWT токен с ролями
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "id":    authUser.Id_user,
        "login": authUser.Login,
        "roles": authUser.Roles,
        "exp":   time.Now().Add(time.Hour * 72).Unix(),
    })

    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
   
    response := Token{
        Token:     tokenString,
        Id_user:   authUser.Id_user,
        Roles:     authUser.Roles,
    }
    type AuthResponse struct {
        Authorized bool     `json:"authorized"`
        UserID     int      `json:"user_id"`
        Login      string   `json:"login"`
        Roles      []string `json:"roles"`
        Message    string   `json:"message"`
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
func VerificationToken(w http.ResponseWriter, r *http.Request) {
    // Нормализация заголовка
    authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
    if authHeader == "" {
        http.Error(w, "Authorization header is required", http.StatusUnauthorized)
        return
    }

    // Гибкое извлечение токена
    var tokenString string
    if strings.HasPrefix(authHeader, "Bearer ") {
        tokenString = strings.TrimPrefix(authHeader, "Bearer ")
    } else {
        tokenString = authHeader
    }

    // Парсинг токена
    token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtKey, nil
    })

    if err != nil {
        http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
        return
    }

    // Проверка валидности claims
    claims, ok := token.Claims.(*jwt.MapClaims)
    if !ok || !token.Valid {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }

    // Извлечение данных
    id, ok := (*claims)["id"].(float64)
    if !ok || id == 0 {
        http.Error(w, "Invalid token: missing user ID", http.StatusUnauthorized)
        return
    }

    login, ok := (*claims)["login"].(string)
    if !ok || login == "" {
        http.Error(w, "Invalid token: missing login", http.StatusUnauthorized)
        return
    }

    // Проверка пользователя в БД
    user, err := db.GetUser(login)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            http.Error(w, "User not found", http.StatusUnauthorized)
            return
        }
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }

    if user.Id_user != int(id) {
        http.Error(w, "Token-user mismatch", http.StatusUnauthorized)
        return
    }

    // Проверка ролей
   
    if roles, ok := (*claims)["roles"].([]interface{}); ok {
        for _, role := range roles {
            if roleMap, ok := role.(map[string]interface{}); ok {
                if name, ok := roleMap["name"].(string); ok && name == "admin" {
                    // Perform some action if the user is an admin
                    http.Error(w, "Admin access granted", http.StatusOK)
                    return
                }
            }
        }
    }

    // Ответ
   // Собираем все роли пользователя
   var roles []string
   if rolesRaw, ok := (*claims)["roles"].([]interface{}); ok {
       for _, r := range rolesRaw {
           if roleMap, ok := r.(map[string]interface{}); ok {
               if roleName, ok := roleMap["name"].(string); ok {
                   // Убираем лишние пробелы и проверяем длину
                   roleName = strings.TrimSpace(roleName)
                   if roleName != "" {
                       roles = append(roles, roleName)
                   }
               }
           }
       }
   }

   // Формируем ответ
   responseText := "Authorized as user" // базовый текст
   if len(roles) > 0 {
       // Упорядочиваем роли: админ всегда первый
       sort.Slice(roles, func(i, j int) bool {
           return roles[i] == "admin" || (roles[i] == "manager" && roles[j] == "user")
       })
       
       // Объединяем уникальные роли
       uniqueRoles := make([]string, 0)
       seen := make(map[string]bool)
       for _, role := range roles {
           if !seen[role] {
               seen[role] = true
               uniqueRoles = append(uniqueRoles, role)
           }
       }
       
       responseText = "Authorized as " + strings.Join(uniqueRoles, ", ")
   }

   w.Header().Set("Content-Type", "text/plain")
   w.Write([]byte(responseText))
}


func PutForgotPassword(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPut {
        http.Error(w, "Invalid request method", http.StatusBadRequest)
        return 
    }
    var authRequest struct {
        Login    string `json:"login"`
        Password string `json:"password"`
    }
    err := json.NewDecoder(r.Body).Decode(&authRequest)
    if err != nil {
        http.Error(w, "Invalid request body PutChangePassword", http.StatusBadRequest)
        return
    }
    if authRequest.Login == "" || authRequest.Password == "" {
        http.Error(w, "Invalid request body PutChangePassword", http.StatusBadRequest)
        return
    }
    user, err := db.GetUser(authRequest.Login)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = db.ChangePassword(user.Login, authRequest.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusBadRequest)
        return
    }
    email := r.URL.Query().Get("email")
    if email == "" {
        http.Error(w, "Email is required", http.StatusBadRequest)
        return
    }
    user, err := db.GetUser(email)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(user)
}