// stages.go
package handlers

import (
	db "appContract/pkg/db/repository"
	"appContract/pkg/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

// Stages
func GetAllStages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method GetAllStages", http.StatusBadRequest)
		return
	}

	stages, err := db.DBgetStageAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var stagesResponse []map[string]interface{}

	for _, stage := range stages {

		stageResponse := map[string]interface{}{
			"id_stage":             stage.Id_stage,
			"name_stage":           stage.Name_stage,
			"id_user":              stage.Id_user,
			"surname":              stage.Surname,
			"username":             stage.Username,
			"patronymic":           stage.Patronymic,
			"phone":                stage.Phone,
			"email":                stage.Email,
			"description":          stage.Description,
			"status_stage":         stage.Id_status_stage,
			"date_change_status":   stage.Date_change_status,
			"name_status_stage":    stage.Name_status_stage,
			"date_create_start":    stage.Date_create_start,
			"date_create_end":      stage.Date_create_end,
			"id_contract":          stage.Id_contract,
			"name_contract":        stage.Name_contract,
			"date_create_contract": stage.Data_contract_create,
		}
		stagesResponse = append(stagesResponse, stageResponse)
	}
	data, ererr := json.Marshal(stagesResponse)
	if ererr != nil {
		http.Error(w, ererr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func GetStagesByIdContract(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	contractID := vars["contractID"]
	if contractID == "" {
		http.Error(w, "Contract ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(contractID)
	if err != nil {
		http.Error(w, "Invalid contract ID format", http.StatusBadRequest)
		return
	}

	stages, err := db.DBgetStageByContractID(id)
	if err != nil {
		log.Printf("Error getting stages: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Подготовка ответа
	response := make([]map[string]interface{}, 0, len(stages))
	for _, stage := range stages {
		response = append(response, map[string]interface{}{
			"id_stage":             stage.Id_stage,
			"name_stage":           stage.Name_stage,
			"id_user":              stage.Id_user,
			"surname":              stage.Surname,
			"username":             stage.Username,
			"patronymic":           stage.Patronymic,
			"description":          stage.Description,
			"id_status_stage":      stage.Id_status_stage,
			"name_status_stage":    stage.Name_status_stage,
			"date_change_status":   stage.Date_change_status,
			"date_create_start":    stage.Date_create_start,
			"date_create_end":      stage.Date_create_end,
			"date_create_contract": stage.Data_contract_create,
			"id_contract":          stage.Id_contract,
			"name_contract":        stage.Name_contract,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func GetUserStages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	userID := vars["userID"]
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid User ID format", http.StatusBadRequest)
		return
	}

	stages, err := db.DBgetStageUserID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем кастомный ответ только с нужными полями
	var responses []map[string]interface{}
	for _, stage := range stages {
		response := map[string]interface{}{
			"id_stage":          stage.Id_stage,
			"name_stage":        stage.Name_stage,
			"description":       stage.Description,
			"date_create_start": stage.Date_create_start,
			"date_create_end":   stage.Date_create_end,
			"name_contract":     stage.Name_contract,
			"name_status_stage": stage.Name_status_stage,
			"surname":           stage.Surname,
			"username":          stage.Username,
			"patronymic":        stage.Patronymic,
		}
		responses = append(responses, response)

	}
	data, err := json.Marshal(responses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func GetStage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method GetStage", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	stageID := vars["stageID"]
	if stageID == "" {
		http.Error(w, "Invalid stage_id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(stageID)
	if err != nil {
		http.Error(w, "Invalid stage_id", http.StatusBadRequest)
		return
	}

	stage, err := db.DBgetStageID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var stageResponse map[string]interface{}

	stageResponse = map[string]interface{}{
		"id_stage":             stage.Id_stage,
		"name_stage":           stage.Name_stage,
		"id_user":              stage.Id_user,
		"surname":              stage.Surname,
		"username":             stage.Username,
		"patronymic":           stage.Patronymic,
		"phone":                stage.Phone,
		"email":                stage.Email,
		"description":          stage.Description,
		"id_status_stage":      stage.Id_status_stage,
		"name_status_stage":    stage.Name_status_stage,
		"date_create_start":    stage.Date_create_start,
		"date_create_end":      stage.Date_create_end,
		"id_contract":          stage.Id_contract,
		"name_contract":        stage.Name_contract,
		"date_create_contract": stage.Data_contract_create,
		"id_type_contract":     stage.Id_type_contract,
		"name_type_contract":   stage.Name_type_contract,
	}

	data, ererr := json.Marshal(stageResponse)
	if ererr != nil {
		http.Error(w, ererr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
func GetStageFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method GetStageFiles", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	stageID := vars["stageID"]
	if stageID == "" {
		http.Error(w, "Invalid stage_id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(stageID)
	if err != nil {
		http.Error(w, "Invalid stage_id", http.StatusBadRequest)
		return
	}
	files, err := db.DBgetFilesStageID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var filesResponse []map[string]interface{}

	for _, file := range files {

		fileResponse := map[string]interface{}{
			"id_file":   file.Id_file,
			"name_file": file.Name_file,
			"data":      file.Data,
			"type_file": file.Type_file,
			"id_stage":  file.Id_stage,
		}
		filesResponse = append(filesResponse, fileResponse)
	}

	data, ererr := json.Marshal(filesResponse)
	if ererr != nil {
		http.Error(w, ererr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func GetStageFilesID(w http.ResponseWriter, r *http.Request) {
	// Проверка метода запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлечение параметров из URL
	vars := mux.Vars(r)
	stageID := vars["stageID"]
	fileID := vars["fileID"]

	// Конвертация параметров в числа
	idStage, err := strconv.Atoi(stageID)
	if err != nil {
		http.Error(w, "Invalid Stage ID format", http.StatusBadRequest)
		return
	}

	idFile, err := strconv.Atoi(fileID)
	if err != nil {
		http.Error(w, "Invalid File ID format", http.StatusBadRequest)
		return
	}

	// Получение файла из БД
	file, err := db.DBgetFileIDStageID(idStage, idFile)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		log.Printf("Database error: %v", err)
		return
	}

	w.Header().Set("Content-Type", file.Type_file)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Name_file))
	w.Header().Set("Content-Length", strconv.Itoa(len(file.Data)))

	// Отправка файла
	if _, err := w.Write(file.Data); err != nil {
		log.Printf("Failed to send file: %v", err)
		http.Error(w, "File transmission failed", http.StatusInternalServerError)
	}
}
func GetStageStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method GetStageStatus", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	statusID := vars["statusID"]
	if statusID == "" {
		http.Error(w, "Invalid status_id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(statusID)
	if err != nil {
		http.Error(w, "Invalid status_id", http.StatusBadRequest)
		return
	}

	status, err := db.DBgetStageIdStatus(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var statusResponse map[string]interface{}

	statusResponse = map[string]interface{}{
		"id_status_stage":   status.Id_status_stage,
		"name_status_stage": status.Name_status_stage,
	}

	data, ererr := json.Marshal(statusResponse)
	if ererr != nil {
		http.Error(w, ererr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method GetStageComments", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	stageID := vars["stageID"]
	if stageID == "" {
		http.Error(w, "Invalid stage_id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(stageID)
	if err != nil {
		log.Panicln(err)
		http.Error(w, "Invalid stage_id", http.StatusBadRequest)
		return
	}
	comments, err := db.DBgetComment(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var commentsResponse []map[string]interface{}
	for _, comment := range comments {
		commentResponse := map[string]interface{}{
			"id_history_state": comment.Id_history_state,
			"id_status_stage":  comment.Id_status_stage,
			"id_stage":         comment.Id_stage,
			"data_create":      comment.Data_create,
			"comment":          comment.Comment,
		}
		commentsResponse = append(commentsResponse, commentResponse)
	}

	data, ererr := json.Marshal(commentsResponse)
	if ererr != nil {
		http.Error(w, ererr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func PostFileToStage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Увеличиваем лимит размера запроса
	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB
		http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Получаем stage_id из URL
	vars := mux.Vars(r)
	stageID, err := strconv.Atoi(vars["stage_id"])
	if err != nil {
		http.Error(w, "Invalid stage ID", http.StatusBadRequest)
		return
	}

	// Получаем файл из запроса
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Читаем ВСЕ данные файла
	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем объект файла
	newFile := models.File{
		Name_file: handler.Filename,
		Data:      fileData,
		Type_file: handler.Header.Get("Content-Type"),
		Id_stage:  stageID,
	}

	// Сохраняем в БД
	if err := db.DBaddFile(newFile); err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "File uploaded successfully",
		"name":    handler.Filename,
		"size":    fmt.Sprintf("%d bytes", handler.Size),
	})
}

func PostCreateStage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method CreateStage", http.StatusBadRequest)
		return
	}
	var stage models.Stages
	err := json.NewDecoder(r.Body).Decode(&stage)
	if err != nil {
		log.Panicln(err)
		http.Error(w, "Invalid request body CreateStage", http.StatusBadRequest)
		return
	}
	err = db.DBaddStage(stage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Stage created successfully"})
}

func PostAddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method CreateComment", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	idStage, err := strconv.Atoi(vars["stageID"])
	if err != nil {
		http.Error(w, "Invalid idStage", http.StatusBadRequest)
		return
	}

	idStatusStage, err := strconv.Atoi(vars["idStatusStage"])
	if err != nil {
		http.Error(w, "Invalid idStatusStage", http.StatusBadRequest)
		return
	}
	var comment models.Stages

	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Invalid request body CreateComment", http.StatusBadRequest)
		return
	}

	err = db.DBaddComment(idStage, idStatusStage, comment.Comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment created successfully"})
}

func PutStageStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request body PutStageStatus", http.StatusBadRequest)
		return
	}
	var stage models.Stages

	err := json.NewDecoder(r.Body).Decode(&stage)
	if err != nil {
		http.Error(w, "Invalid request body PutStageStatus", http.StatusBadRequest)
		return
	}

	if stage.Id_stage == 0 || stage.Id_status_stage == 0 || stage.Comment == "" {
		http.Error(w, "Invalid request body PutStageStatus", http.StatusBadRequest)
		return
	}

	err = db.DBChengeStatusStage(stage.Id_stage, stage.Id_status_stage, stage.Comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Stage status updated successfully"})
}

func DeleteStageFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method DeleteStageFiles", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	id_faile, err := strconv.Atoi(vars["id_file"])
	if err != nil {
		http.Error(w, "Invalid id_file", http.StatusBadRequest)
		return
	}
	err = db.DBdeleteFile(id_faile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "File deleted successfully"})

}
func DeleteStage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method DeleteStage", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	stageID, err := strconv.Atoi(vars["stageID"])
	if err != nil {
		http.Error(w, "Invalid stage_id", http.StatusBadRequest)
		return
	}
	err = db.DBdeleteStage(stageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Stage deleted successfully"})
	// Дополнительная проверка: если stageID не существует, вернуть ошибку
	if stageID == 0 {
		http.Error(w, "Stage not found", http.StatusNotFound)
		return
	}
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method DeleteComment", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	idComment, err := strconv.Atoi(vars["idComment"])
	if err != nil {
		http.Error(w, "Invalid idComment", http.StatusBadRequest)
		return
	}
	err = db.DBdeleteComment(idComment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment deleted successfully"})
}
