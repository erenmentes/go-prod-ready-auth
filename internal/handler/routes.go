package handler

import (
	"github.com/erenmentes/go-prod-ready-auth/internal/middleware"
	"github.com/erenmentes/go-prod-ready-auth/internal/service"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetupRoutes(mux *mux.Router, db *gorm.DB) error {
	mailService := service.NewMailService()

	authService := service.NewAuthService(mailService, db)

	authHandler := NewAuthHandler(authService)
	authMiddleware := middleware.NewAuthMiddleware(db)
	emailMiddleware := middleware.NewEmailMiddleware()

	// Auth Endpoints
	mux.HandleFunc("/login", authHandler.Login).Methods("POST")
	mux.HandleFunc("/login/2fa", authHandler.VerifyTwoFactor).Methods("POST")
	mux.HandleFunc("/register", authHandler.Register).Methods("POST")
	mux.HandleFunc("/refresh-token", authHandler.RefreshToken).Methods("POST")
	mux.HandleFunc("/verify-account", authHandler.VerifyAccount).Methods("GET")
	mux.HandleFunc("/resend-verification-email", authHandler.ResendVerificationEmail).Methods("POST")
	mux.HandleFunc("/toggle-2fa", authMiddleware.Apply(emailMiddleware.Apply(authHandler.ToggleTwoFactor))).Methods("POST")
	mux.HandleFunc("/reset-password", authMiddleware.Apply(emailMiddleware.Apply(authHandler.ResetPassword))).Methods("POST")

	return nil
}
