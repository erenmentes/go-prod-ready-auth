package service

type MailService struct{}

func NewMailService() *MailService {
	return &MailService{}
}

func (s *MailService) SendAccountCreationVerificationCode(to string) error {
	return nil
}

func (s *MailService) SendTwoFactorVerificationMail(to string) error {
	return nil
}
