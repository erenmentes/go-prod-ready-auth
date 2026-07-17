package handler

import (
	"github.com/erenmentes/go-prod-ready-auth/internal/service"
	"github.com/gorilla/mux"
)

func SetupRoutes(mux *mux.Router) error {
	mailService := service.NewMailService()

	authService := service.NewAuthService(mailService)

	authHandler := NewAuthHandler(authService)

	// Auth Endpoints
	mux.HandleFunc("/login", authHandler.Login).Methods("POST")
	mux.HandleFunc("/register", authHandler.Register).Methods("POST")
	mux.HandleFunc("/refresh-token", authHandler.RefreshToken).Methods("POST")

	return nil
}
