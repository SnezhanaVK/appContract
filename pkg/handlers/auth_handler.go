package handlers

import (
	db "appContract/pkg/db/repository"
	"appContract/pkg/models"
	"appContract/pkg/service"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Обновленный обработчик Login
var jwtKey = []byte("secretkey")

// AuthResponse структура для ответа после авторизации (без токена)
type AuthResponse struct {
	Id_user int  `json:"id_user"`
	Admin   bool `json:"admin"`
	Manager bool `json:"manager"`
}

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

	// Создаем JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      authUser.Id_user,
		"login":   authUser.Login,
		"admin":   authUser.Admin,
		"manager": authUser.Manager,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем токен в куки
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(72 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // В production должно быть true
		SameSite: http.SameSiteLaxMode,
	})

	// Отправляем ответ с остальными данными пользователя
	response := AuthResponse{
		Id_user: authUser.Id_user,
		Admin:   authUser.Admin,
		Manager: authUser.Manager,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func VerificationToken(w http.ResponseWriter, r *http.Request) {
    // Получаем токен из куки
    cookie, err := r.Cookie("token")
    if err != nil {
        if err == http.ErrNoCookie {
            respondWithBool(w, false)
            return
        }
        respondWithBool(w, false)
        return
    }

    tokenString := cookie.Value

    // Парсинг токена
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtKey, nil
    })

    if err != nil || !token.Valid {
        respondWithBool(w, false)
        return
    }

    // Если токен валиден
    respondWithBool(w, true)
}

func respondWithBool(w http.ResponseWriter, result bool) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]bool{"valid": result})
}
func Logout(w http.ResponseWriter, r *http.Request) {
	// Создаем куку с таким же именем ("token"), но с истекшим сроком действия
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),  // Дата в прошлом (мгновенно истекает)
		MaxAge:   -1,               // Немедленное удаление куки
		HttpOnly: true,
		Secure:   true,             // Должно совпадать с настройками при логине
		SameSite: http.SameSiteLaxMode,
	})

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func PutForgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	var authRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if authRequest.Email == "" || authRequest.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Проверяем существование пользователя
	_, err := db.GetUser(authRequest.Email)
	if err != nil {
		if err.Error() == "User not found" {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Передаем email из запроса, а не user.Login!
	if err := db.ChangePassword(authRequest.Email, authRequest.Password); err != nil {
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

func PostSendEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}
	var emailRequest struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&emailRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	var user models.Users
	user, err = db.GetUserByEmail(emailRequest.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	service.SendingCode(user)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Email sent successfully"})
}

func PostVerifyCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	var request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if service.VerifyCode(request.Email, request.Code) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Code verified successfully"})
	} else {
		http.Error(w, "Invalid or expired code", http.StatusUnauthorized)
	}
}
