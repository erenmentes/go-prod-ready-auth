package service

import (
	"errors"
	"time"

	"github.com/erenmentes/go-prod-ready-auth/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IMailService interface {
	SendAccountCreationVerificationCode(to, emailVerificationCode string) error
	SendTwoFactorVerificationMail(to string) error
}

type AuthService struct {
	IMailService
	db *gorm.DB
}

func NewAuthService(mailService IMailService, db *gorm.DB) *AuthService {
	return &AuthService{
		IMailService: mailService,
		db:           db,
	}
}

func (s *AuthService) Login(username, password string) error {
	return nil
}

func (s *AuthService) Register(email, username, password string) error {

	var EmailVerificationCode uuid.UUID

	exists, err := s.checkIfUserAlreadyExists(email)
	if err != nil {
		return errors.New("Something went wrong during user availability check.")
	}
	if exists {
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
		return errors.New("Something went wrong while creating user in database.")
	}

	if err := s.SendAccountCreationVerificationCode(email, EmailVerificationCode.String()); err != nil {
		tx.Rollback()
		return errors.New("Something went wrong while sending verification code to user.")
	}

	tx.Commit()

	return nil
}

func (s *AuthService) RefreshToken(refreshToken string) error {
	return nil
}

func (s *AuthService) VerifyAccount(verificationCode string) error {
	return nil
}

func (s *AuthService) VerifyTwoFactorVerification(verificationCode string) error {
	return nil
}

func (s *AuthService) checkIfUserAlreadyExists(email string) (bool, error) {
	var exists bool

	err := s.db.Model(&models.User{}).
		Select("count(1) > 0").
		Where("email = ?", email).
		Find(&exists).Error

	if err != nil {
		return false, err
	}

	return exists, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
