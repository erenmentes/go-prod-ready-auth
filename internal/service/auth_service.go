package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/erenmentes/go-prod-ready-auth/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IMailService interface {
	SendAccountCreationVerificationCode(to, emailVerificationCode string) error
	SendTwoFactorVerificationMail(to string) (string, error)
}

type AuthService struct {
	IMailService
	db *gorm.DB
}

type JwtPayload struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

const TOKEN_DURATION = 15 * time.Minute

var JWT_SECRET = os.Getenv("JWT_SECRET")

func NewAuthService(mailService IMailService, db *gorm.DB) *AuthService {
	return &AuthService{
		IMailService: mailService,
		db:           db,
	}
}

func (s *AuthService) Login(email, password string) (*LoginResponse, error) {
	user, err := s.getUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if user == nil {
		// fake hash with valid format 🦫, they ain't gettin that email datas from our backend bud!!!
		fakeHash := "$2a$10$AzsX4.zN2v7YgqXvO1wGDu6V6p6l6f6v6f6v6f6v6f6v6f6v6f6v6"
		_ = bcrypt.CompareHashAndPassword([]byte(fakeHash), []byte(password))
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if user.IsTwoFactorVerificationActivated {
		verificationCode, err := s.SendTwoFactorVerificationMail(user.Email)
		if err != nil {
			return nil, errors.New("failed to send two-factor verification code")
		}

		expiry := time.Now().Add(15 * time.Minute)
		err = s.db.Model(&models.User{}).
			Where("id = ?", user.ID).
			Updates(models.User{
				TwoFactorVerificationCode:       &verificationCode,
				TwoFactorVerificationExpiryDate: &expiry,
			}).Error
		if err != nil {
			return nil, errors.New("failed to persist two-factor verification code")
		}

		return nil, errors.New("two factor verification required")
	}

	claims := JwtPayload{
		UserID: uint(user.ID),
		Email:  user.Email,
		Role:   user.UserRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_DURATION)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := accessToken.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken := uuid.New().String()

	err = s.db.Model(&models.User{}).
		Where("id = ?", user.ID).
		Update("refreshtoken", refreshToken).Error
	if err != nil {
		return nil, errors.New("failed to persist session")
	}

	return &LoginResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Register(email, username, password string) error {

	var EmailVerificationCode uuid.UUID

	user, err := s.getUserByEmail(email)
	if err != nil {
		return errors.New("Something went wrong during user check.")
	}
	if user != nil {
		return errors.New("User with this email already exists.")
	}

	hashedPass, err := HashPassword(password)

	if err != nil {
		return errors.New("Something went wrong while hashing password.")
	}

	EmailVerificationCode = uuid.New()

	newUser := models.User{
		Username:                        username,
		Email:                           email,
		Pass:                            hashedPass,
		EmailVerificationCode:           EmailVerificationCode.String(),
		EmailVerificationCodeExpiryDate: time.Now().Add(24 * time.Hour),
	}

	tx := s.db.Begin()

	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		return errors.New("something went wrong while creating user in database.")
	}

	if err := s.SendAccountCreationVerificationCode(email, EmailVerificationCode.String()); err != nil {
		tx.Rollback()
		return errors.New("Something went wrong while sending verification code to user.")
	}

	tx.Commit()

	return nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*RefreshTokenResponse, error) {
	var user models.User

	err := s.db.Where("RefreshToken = ?", refreshToken).First(&user).Error
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	newRefreshToken := uuid.New().String()

	err = s.db.Model(&models.User{}).
		Where("id = ?", user.ID).
		Update("RefreshToken", newRefreshToken).Error
	if err != nil {
		return nil, errors.New("something went wrong while updating session")
	}

	claims := JwtPayload{
		UserID: uint(user.ID),
		Email:  user.Email,
		Role:   user.UserRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_DURATION)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := accessToken.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	return &RefreshTokenResponse{
		AccessToken:  tokenString,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) ToggleTwoFactorVerification(userID uint, activated bool) error {
	result := s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("istwofactorverificationactivated", activated)
	if result.Error != nil {
		return errors.New("failed to update two-factor verification setting")
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (s *AuthService) VerifyAccount(verificationCode string) error {
	var user models.User

	err := s.db.Where("emailverificationcode = ? AND emailverificationcodeexpirydate >= ?", verificationCode, time.Now()).First(&user).Error
	if err != nil {
		return errors.New("invalid or expired verification code")
	}

	verified := true
	user.IsEmailVerified = &verified

	newCode := uuid.New().String()
	err = s.db.Model(&user).
		Select("IsEmailVerified", "EmailVerificationCode").
		Updates(models.User{
			IsEmailVerified:       &verified,
			EmailVerificationCode: newCode,
		}).Error

	if err != nil {
		return errors.New("something went wrong while verifying account")
	}

	return nil
}

func (s *AuthService) VerifyTwoFactorVerification(verificationCode string) (*LoginResponse, error) {
	var user models.User

	err := s.db.Where("twofactorverificationcode = ? AND twofactorverificationexpirydate >= ? AND istwofactorverificationactivated = ?", verificationCode, time.Now(), true).First(&user).Error
	if err != nil {
		return nil, errors.New("invalid or expired two-factor verification code")
	}

	accessToken, err := s.generateAccessToken(&user)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken := uuid.New().String()

	err = s.db.Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(models.User{
			RefreshToken:              refreshToken,
			TwoFactorVerificationCode: nil,
		}).Error
	if err != nil {
		return nil, errors.New("failed to persist session")
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) generateAccessToken(user *models.User) (string, error) {
	claims := JwtPayload{
		UserID: uint(user.ID),
		Email:  user.Email,
		Role:   user.UserRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_DURATION)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWT_SECRET))
}

func (s *AuthService) ResendAccountVerificationEmail(email string) error {
	user, err := s.getUserByEmail(email)
	if err != nil {
		return errors.New("something went wrong while retrieving user")
	}

	if user == nil {
		return errors.New("user not found")
	}

	if user.IsEmailVerified != nil && *user.IsEmailVerified {
		return errors.New("account already verified")
	}

	newCode := uuid.New().String()
	expiry := time.Now().Add(24 * time.Hour)

	tx := s.db.Begin()
	if tx.Error != nil {
		return errors.New("failed to start database transaction")
	}

	err = tx.Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(models.User{
			EmailVerificationCode:           newCode,
			EmailVerificationCodeExpiryDate: expiry,
		}).Error
	if err != nil {
		tx.Rollback()
		return errors.New("something went wrong while updating verification code")
	}

	err = s.SendAccountCreationVerificationCode(user.Email, newCode)
	if err != nil {
		tx.Rollback()
		return errors.New("something went wrong while sending verification code to user")
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("failed to finalize verification code update")
	}

	return nil
}

func (s *AuthService) ResetPassword(email, currentPassword, newPassword, newPasswordAgain string) error {

	user, err := s.getUserByEmail(email)
	if err != nil {
		return errors.New("something went wrong while retrieving user")
	}

	if user == nil {
		return errors.New("user not found")
	}

	e := bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(currentPassword))

	if e != nil {
		return errors.New("invalid current password")
	}

	if newPassword != newPasswordAgain {
		return errors.New("passwords does not match")
	}

	hashedNewPass, err := HashPassword(newPassword)

	if err != nil {
		return errors.New("something went wrong while hashing password")
	}

	dbErr := s.db.Model(&models.User{}).Where("email = ?", email).Updates(models.User{
		Pass: hashedNewPass,
	}).Error

	if dbErr != nil {
		return errors.New("something went wrong while resetting password")
	}

	return nil
}

func (s *AuthService) getUserByEmail(email string) (*models.User, error) {
	var user models.User

	err := s.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
