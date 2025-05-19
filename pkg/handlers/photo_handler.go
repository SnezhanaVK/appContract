package handlers

import (
	db "appContract/pkg/db/repository"
	"appContract/pkg/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetPhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	photo, err := db.DBgetPhoto(userID)
	if err != nil {
		http.Error(w, "Failed to get photo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if photo == nil {
		http.Error(w, "Avatar not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", photo.Type_photo)
	w.Header().Set("Content-Length", strconv.Itoa(len(photo.Data_photo)))

	if _, err := w.Write(photo.Data_photo); err != nil {
		log.Printf("Failed to write image: %v", err)
	}
}
func PostAddPhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid  ID user", http.StatusBadRequest)
		return
	}

	photo, handler, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Failed to get photo: "+err.Error(), http.StatusBadRequest)
		return
	}
	photoData, err := io.ReadAll(photo)
	if err != nil {
		http.Error(w, "Failed to read photo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer photo.Close()

	photoModel := models.Photo{
		Data_photo: photoData,
		Type_photo: handler.Header.Get("Content-Type"),
		Id_user:    userID,
	}

	err = db.DBChangePhoto(photoModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Photo added successfully",
		"name":    handler.Filename,
		"size":    strconv.FormatInt(handler.Size, 10),
		"type":    handler.Header.Get("Content-Type"),
	})
}

func PutChangePhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid  ID user", http.StatusBadRequest)
		return
	}

	photo, handler, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Failed to get photo: "+err.Error(), http.StatusBadRequest)
		return
	}

	photoData, err := io.ReadAll(photo)
	if err != nil {
		http.Error(w, "Failed to read photo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer photo.Close()

	photoModel := models.Photo{
		Data_photo: photoData,
		Type_photo: handler.Header.Get("Content-Type"),
		Id_user:    userID,
	}

	err = db.DBChangePhoto(photoModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Photo changed successfully",
		"name":    handler.Filename,
		"size":    strconv.FormatInt(handler.Size, 10),
		"type":    handler.Header.Get("Content-Type"),
	})
}

func DeletePhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = db.DBDeletePhoto(userID)
	if err != nil {
		http.Error(w, "Failed to delete photo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Photo deleted successfully",
	})
}
