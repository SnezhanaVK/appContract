package handlers

// notifications_test_handler.go в папке handlers

import (
	service "appContract/pkg/service"
	"fmt"
	"net/http"
)




func TestNotifications(w http.ResponseWriter, r *http.Request) {
	if err := service.ProcessDailyNotifications(); err != nil { // Теперь без параметра
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	fmt.Fprint(w, "Notifications processed successfully")
}