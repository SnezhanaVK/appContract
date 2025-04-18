package handlers

import (
	service "appContract/pkg/services"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)
func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the appContract server!"))
}
var upgrader = websocket.Upgrader{} // объявляем переменную upgrader
func Search(w http.ResponseWriter, r *http.Request) {
    // Получить данные из запроса
    var searchCriteria struct {
        NameContract string `json:"name_contract"`
        NameStage    string `json:"name_stage"`
        NameTeg      string `json:"name_teg"`
    }
    err := json.NewDecoder(r.Body).Decode(&searchCriteria)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Создать WebSocket соединение
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer ws.Close()

    // Обрабатывать сообщения от клиента
    for {
        messageType, p, err := ws.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }

        // Обработать сообщение от клиента
        if messageType == websocket.TextMessage {
            var message struct {
                SearchCriteria struct {
                    NameContract string `json:"name_contract"`
                    NameStage    string `json:"name_stage"`
                    NameTeg      string `json:"name_teg"`
                } `json:"search_criteria"`
            }
            err := json.Unmarshal(p, &message)
            if err != nil {
                log.Println(err)
                return
            }

            // Выполнить поиск
            results := service.SearchContract(message.SearchCriteria.NameContract, message.SearchCriteria.NameStage, message.SearchCriteria.NameTeg)

            // Отправить результаты обратно клиенту
            data, err := json.Marshal(results)
            if err != nil {
                log.Println(err)
                return
            }
            err = ws.WriteMessage(websocket.TextMessage, data)
            if err != nil {
                log.Println(err)
                return
            }
        }
    }
}
