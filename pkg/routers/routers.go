package routers

import (
	"appContract/pkg/handlers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
    router := mux.NewRouter()

    // Авторизация
    router.HandleFunc("/api/authorizations", handlers.Login).Methods("POST")
    router.HandleFunc("/api/authorizations/token", handlers.VerificationToken).Methods("GET")
    router.HandleFunc("/api/authorizations/forgot-password", handlers.PutForgotPassword).Methods("PUT")
    router.HandleFunc("/api/authorizations/sendingCode", handlers.PostSendEmail).Methods("POST")
    router.HandleFunc("/api/authorizations/verifyCode", handlers.PostVerifyCode).Methods("POST")
    

    // Вызов маршрутов
    router.HandleFunc("/", handlers.Index).Methods("GET")
    router.HandleFunc("api/search", handlers.Search).Methods("POST")

    UserRoutes(router)
    ContractRoutes(router)
    StageRoutes(router)
    ContractsAndStageRoutes(router)

    // Добавьте новый роут
// Тестовый эндпоинт

	// Тестовый эндпоинт
    router.HandleFunc("/test-notifications", handlers.TestNotifications).Methods("GET")

    

    return router
}