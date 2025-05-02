package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	db "appContract/pkg/db/repository"

	"github.com/gorilla/mux"
)

type NotificationRequest struct {
	Days []int `json:"days"`
}

func GetSettings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	settings, err := db.GetUserNotificationSettings(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":  userID,
		"settings": settings,
	})
}

func PutUpdateUserSettings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for _, day := range req.Days {
		if day != 1 && day != 3 && day != 7 {
			http.Error(w, "Invalid notification day. Allowed values: 1, 3, 7", http.StatusBadRequest)
			return
		}
	}

	if err := db.SetUserNotificationSettings(userID, req.Days); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Notification settings updated successfully")
}

func DeleteSettings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := db.SetUserNotificationSettings(userID, []int{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "All notification settings removed")
}
