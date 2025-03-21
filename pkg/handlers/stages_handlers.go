// stages.go
package handlers

import (
	db "appContract/pkg/db/repository"
	"appContract/pkg/models"
	"encoding/json"
	"net/http"
	"strconv"
)

// Stages
func GetAllStages(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method GetAllStages", http.StatusBadRequest)
        return
    }

    stages, err:= db.DBgetStageAll()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
        w.WriteHeader(http.StatusOK)
        if err := json.NewEncoder(w).Encode(stages); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
}

func GetUserStages(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet{
        http.Error(w,"Invalid request method GetUserStages",http.StatusBadRequest)
        return
    }
    userID, err:=strconv.Atoi(r.URL.Query().Get("user_id"))
    if err != nil {
        http.Error(w, "Invalid user_id", http.StatusBadRequest)
        return
    }
    stages, err := db.DBgetStageUserID(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(stages); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func GetStage(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method GetStage", http.StatusBadRequest)
        return
    }
    stageID, err:=strconv.Atoi(r.URL.Query().Get("stage_id"))
    if err != nil {
        http.Error(w, "Invalid stage_id", http.StatusBadRequest)
        return
    }
    stage, err := db.DBgetStageID(stageID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(stage)
}

func GetStageFiles(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method GetStageFiles", http.StatusBadRequest)
        return
    }
    stageID, err:=strconv.Atoi(r.URL.Query().Get("stage_id"))
    if err != nil {
        http.Error(w, "Invalid stage_id", http.StatusBadRequest)
        return
    }
    files, err := db.DBgetFilesStageID(stageID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(files); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func GetStageStatus(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method GetStageStatus", http.StatusBadRequest)
        return
    }
    stageID, err:=strconv.Atoi(r.URL.Query().Get("stage_id"))
    if err != nil {
        http.Error(w, "Invalid stage_id", http.StatusBadRequest)
        return
    }
    status, err := db.DBgetStageIdStatus(stageID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(status); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func GetComments(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method GetStageComments", http.StatusBadRequest)
        return
    }
    stageID, err:=strconv.Atoi(r.URL.Query().Get("stage_id"))
    if err != nil {
        http.Error(w, "Invalid stage_id", http.StatusBadRequest)
        return
    }
    comments, err := db.DBgetComment(stageID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(comments); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func PostFileToStage(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method PostFileToStage", http.StatusBadRequest)
        return
    }
    var file models.File
    err := json.NewDecoder(r.Body).Decode(&file)
    if err != nil {
        http.Error(w, "Invalid request body PostFileToStage", http.StatusBadRequest)
        return
    }
    err = db.DBaddFile(file)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "File uploaded successfully"})
}

func PostCreateStage(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method CreateStage", http.StatusBadRequest)
        return
    }
    var stage models.Stages
    err := json.NewDecoder(r.Body).Decode(&stage)
    if err != nil {
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

func PostCreateComment (w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method CreateComment", http.StatusBadRequest)
        return
    }
    var comment models.Stages
    err := json.NewDecoder(r.Body).Decode(&comment)
    if err != nil {
        http.Error(w, "Invalid request body CreateComment", http.StatusBadRequest)
        return
    }
    err = db.DBCreateComment(comment)
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
    err = db.DBChengeStatusStage(stage.Id_stage, stage.Id_status_stage,stage.Comment )
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
    id_faile, err:=strconv.Atoi(r.URL.Query().Get("id_file"))
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
    stageID, err:=strconv.Atoi(r.URL.Query().Get("stage_id"))
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
}




