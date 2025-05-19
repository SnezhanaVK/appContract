package routers

import (
	"appContract/pkg/handlers"
	"appContract/pkg/middleware"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.MonitoringMiddleware)

	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	router.HandleFunc("/api/authorizations", handlers.Login).Methods("POST")
	router.HandleFunc("/api/authorizations/token", handlers.VerificationToken).Methods("GET")
	router.HandleFunc("/api/authorizations/logout", handlers.Logout).Methods("GET")
	router.HandleFunc("/api/authorizations/forgot-password", handlers.PutForgotPassword).Methods("PUT")
	router.HandleFunc("/api/authorizations/sendingCode", handlers.PostSendEmail).Methods("POST")
	router.HandleFunc("/api/authorizations/verifyCode", handlers.PostVerifyCode).Methods("POST")

	router.HandleFunc("/", handlers.Index).Methods("GET")
	router.HandleFunc("/api/search", handlers.Search).Methods("POST")

	UserRoutes(router)
	ContractRoutes(router)
	StageRoutes(router)
	ContractsAndStageRoutes(router)
	SortRouters(router)
	NotificationRoutes(router)
	PhotoRouters(router)

	router.HandleFunc("/test-notifications", handlers.TestNotifications).Methods("GET")

	return router
}
