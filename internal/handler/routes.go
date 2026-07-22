package handler

import (
	"github.com/erenmentes/go-prod-ready-auth/internal/service"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetupRoutes(mux *mux.Router, db *gorm.DB) error {
	mailService := service.NewMailService()

	authService := service.NewAuthService(mailService, db)

	authHandler := NewAuthHandler(authService)

	// Auth Endpoints
	mux.HandleFunc("/login", authHandler.Login).Methods("POST")
	mux.HandleFunc("/register", authHandler.Register).Methods("POST")
	mux.HandleFunc("/refresh-token", authHandler.RefreshToken).Methods("POST")

	return nil
}
