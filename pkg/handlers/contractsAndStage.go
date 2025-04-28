package handlers

// import (
// 	db "appContract/pkg/db/repository"
// 	"appContract/pkg/models"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// )

// func GetContractsAndStages(w http.ResponseWriter, r *http.Request) {
//     if r.Method != http.MethodGet {
//         http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//         return
//     }

//     // Получаем все контракты
//     contracts, err := db.DBgetContractAll()
//     if err != nil {
//         http.Error(w, fmt.Sprintf("Failed to get contracts: %v", err), http.StatusInternalServerError)
//         return
//     }

//     // Получаем все этапы
//     stages, err := db.DBgetStageAll()
//     if err != nil {
//         http.Error(w, fmt.Sprintf("Failed to get stages: %v", err), http.StatusInternalServerError)
//         return
//     }

//     // Создаем мапу для быстрого доступа к этапам по ID контракта
//     stagesByContract := make(map[int][]models.Stages)
//     for _, stage := range stages {
//         stagesByContract[stage.Id_contract] = append(stagesByContract[stage.Id_contract], stage)
//     }

//     var response []map[string]interface{}

//     // Формируем ответ с объединением контрактов и их этапов
//     for _, contract := range contracts {
//         // Подготавливаем теги контракта
//         var tags []map[string]interface{}
//         for _, tag := range contract.Tegs {
//             tags = append(tags, map[string]interface{}{
//                 "id":   tag.Id_tegs,
//                 "name": tag.Name_tegs,
//             })
//         }

//         // Подготавливаем этапы для этого контракта
//         var contractStages []map[string]interface{}
//         for _, stage := range stagesByContract[contract.Id_contract] {
//             stageData := map[string]interface{}{
//                 "id":                stage.Id_stage,
//                 "name":              stage.Name_stage,
//                 "user": map[string]interface{}{
//                     "id":         stage.Id_user,
//                     "surname":    stage.Surname,
//                     "name":       stage.Username,
//                     "patronymic": stage.Patronymic,
//                     "phone":     stage.Phone,
//                     "email":      stage.Email,
//                 },
//                 "description":      stage.Description,
//                 "status": map[string]interface{}{
//                     "id":   stage.Id_status_stage,
//                     "name": stage.Name_status_stage,
//                 },
//                 "dates": map[string]interface{}{
//                     "change_status": stage.Date_change_status,
//                     "start":         stage.Date_create_start,
//                     "end":           stage.Date_create_end,
//                 },
//             }
//             contractStages = append(contractStages, stageData)
//         }

//         // Формируем данные контракта
//         contractData := map[string]interface{}{
//             "id":          contract.Id_contract,
//             "name":        contract.Name_contract,
//             "dates": map[string]interface{}{
//                 "create":      contract.Date_contract_create,
//                 "conclusion":  contract.Date_conclusion,
//                 "end":        contract.Date_end,
//             },
//             "user": map[string]interface{}{
//                 "id":         contract.Id_user,
//                 "surname":    contract.Surname,
//                 "name":       contract.Username,
//                 "patronymic": contract.Patronymic,
//             },
//             "type": map[string]interface{}{
//                 "id":   contract.Id_type,
//                 "name": contract.Name_type,
//             },
//             "counterparty": map[string]interface{}{
//                 "id":   contract.Id_counterparty,
//                 "name": contract.Name_counterparty,
//             },
//             "status": map[string]interface{}{
//                 "id":   contract.Id_status_contract,
//                 "name": contract.Name_status_contract,
//             },
//             "tags":    tags,
//             "stages":  contractStages,
//         }

//         response = append(response, contractData)
//     }

//     w.Header().Set("Content-Type", "application/json")
//     if err := json.NewEncoder(w).Encode(response); err != nil {
//         http.Error(w, "Failed to encode response", http.StatusInternalServerError)
//     }
// }