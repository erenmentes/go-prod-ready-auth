package service

type IMailService interface {
	SendAccountCreationVerificationCode(to string) error
	SendTwoFactorVerificationMail(to string) error
}

type AuthService struct {
	IMailService
}

func NewAuthService(mailService IMailService) *AuthService {
	return &AuthService{
		IMailService: mailService,
	}
}

func (s *AuthService) Login(username, password string) error {
	return nil
}

func (s *AuthService) Register(email, username, password string) error {
	return nil
}

func (s *AuthService) RefreshToken(refreshToken string) error {
	return nil
}
