package handler

import (
	"encoding/json"
	"net/http"

	"github.com/erenmentes/go-prod-ready-auth/internal/middleware"
	"github.com/erenmentes/go-prod-ready-auth/internal/service"
)

type IAuthService interface {
	service.IMailService
	Login(email, password string) (*service.LoginResponse, error)
	Register(email, username, password string) error
	RefreshToken(refreshToken string) (*service.RefreshTokenResponse, error)
	VerifyAccount(verificationCode string) error
	ResendAccountVerificationEmail(email string) error
	VerifyTwoFactorVerification(verificationCode string) (*service.LoginResponse, error)
	ToggleTwoFactorVerification(userID uint, activated bool) error
	ResetPassword(email, currentPassword, newPassword, newPasswordAgain string) error
}

type AuthHandler struct {
	IAuthService
}

func NewAuthHandler(authService IAuthService) *AuthHandler {
	return &AuthHandler{
		IAuthService: authService,
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type verificationRequest struct {
	VerificationCode string `json:"verification_code"`
}

type resendVerificationRequest struct {
	Email string `json:"email"`
}

type toggleTwoFactorRequest struct {
	Enabled bool `json:"enabled"`
}

type resetPasswordRequest struct {
	CurrentPassword  string `json:"current_password"`
	NewPassword      string `json:"new_password"`
	NewPasswordAgain string `json:"new_password_again"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	resp, err := h.IAuthService.Login(req.Email, req.Password)
	if err != nil {
		if err.Error() == "two factor verification required" {
			respondJSON(w, http.StatusAccepted, map[string]string{"message": "two factor verification required"})
			return
		}
		// will also check if the account is locked (brute force protection)
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.IAuthService.Register(req.Email, req.Username, req.Password); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]string{"message": "user registered successfully"})
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req refreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	resp, err := h.IAuthService.RefreshToken(req.RefreshToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) ToggleTwoFactor(w http.ResponseWriter, r *http.Request) {
	var req toggleTwoFactorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthenticated")
		return
	}

	if err := h.IAuthService.ToggleTwoFactorVerification(uint(user.ID), req.Enabled); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "two factor verification updated"})
}

func (h *AuthHandler) VerifyAccount(w http.ResponseWriter, r *http.Request) {
	verificationCode := r.URL.Query().Get("code")
	if verificationCode == "" {
		respondError(w, http.StatusBadRequest, "verification code is required")
		return
	}

	if err := h.IAuthService.VerifyAccount(verificationCode); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "account verified successfully"})
}

func (h *AuthHandler) ResendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	var req resendVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.IAuthService.ResendAccountVerificationEmail(req.Email); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "verification email resent"})
}

func (h *AuthHandler) VerifyTwoFactor(w http.ResponseWriter, r *http.Request) {
	var req verificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	resp, err := h.IAuthService.VerifyTwoFactorVerification(req.VerificationCode)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req resetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthenticated")
		return
	}

	err := h.IAuthService.ResetPassword(user.Email, req.CurrentPassword, req.NewPassword, req.NewPasswordAgain)

	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
