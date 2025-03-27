package routers

import (
	"appContract/pkg/handlers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
    router := mux.NewRouter()

    // Авторизация
    router.HandleFunc("/api/authorizations ", handlers.Login).Methods("POST")
    router.HandleFunc("/api/forgot-password", handlers.PutForgotPassword).Methods("PUT")

    // Вызов маршрутов
    UserRoutes(router)
    ContractRoutes(router)
    StageRoutes(router)
    // NotificationRoutes(router)

    return router
}
