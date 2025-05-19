package handlers

import (
	service "appContract/pkg/service"
	"fmt"
	"net/http"
)

func TestNotifications(w http.ResponseWriter, r *http.Request) {
	if err := service.ProcessDailyNotifications(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	fmt.Fprint(w, "Notifications processed successfully")
}
