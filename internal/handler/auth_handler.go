package handler

import (
	"net/http"

	"github.com/erenmentes/go-prod-ready-auth/internal/service"
)

type IAuthService interface {
	service.IMailService
	Login(email, password string) (*service.LoginResponse, error)
	Register(email, username, password string) error
	RefreshToken(refreshToken string) error
}

type AuthHandler struct {
	IAuthService
}

func NewAuthHandler(authService IAuthService) *AuthHandler {
	return &AuthHandler{
		IAuthService: authService,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {

}
