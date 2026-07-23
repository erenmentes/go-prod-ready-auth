package service

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type MailService struct{}

func NewMailService() *MailService {
	return &MailService{}
}

func (s *MailService) SendAccountCreationVerificationCode(to, emailVerificationCode string) error {
	mailtrapAPIURL := strings.TrimSpace(os.Getenv("MAILTRAP_API_URL"))
	mailtrapToken := strings.TrimSpace(os.Getenv("MAILTRAP_API_TOKEN"))
	mailtrapFromEmail := strings.TrimSpace(os.Getenv("MAILTRAP_FROM_EMAIL"))
	mailtrapFromName := strings.TrimSpace(os.Getenv("MAILTRAP_FROM_NAME"))

	if mailtrapAPIURL == "" || mailtrapToken == "" || mailtrapFromEmail == "" {
		return fmt.Errorf("Mailtrap API configuration is incomplete.")
	}

	payload := map[string]interface{}{
		"from": map[string]string{
			"email": mailtrapFromEmail,
			"name":  mailtrapFromName,
		},
		"to":       []map[string]string{{"email": to}},
		"subject":  "Verify your account",
		"html":     buildVerificationEmailHTML(emailVerificationCode),
		"category": "Account Verification",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, mailtrapAPIURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+mailtrapToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("mailtrap request failed with status: %s", res.Status)
	}

	return nil
}

func (s *MailService) SendTwoFactorVerificationMail(to string) (string, error) {
	mailtrapAPIURL := strings.TrimSpace(os.Getenv("MAILTRAP_API_URL"))
	mailtrapToken := strings.TrimSpace(os.Getenv("MAILTRAP_API_TOKEN"))
	mailtrapFromEmail := strings.TrimSpace(os.Getenv("MAILTRAP_FROM_EMAIL"))
	mailtrapFromName := strings.TrimSpace(os.Getenv("MAILTRAP_FROM_NAME"))

	if mailtrapAPIURL == "" || mailtrapToken == "" || mailtrapFromEmail == "" {
		return "", fmt.Errorf("Mailtrap API configuration is incomplete.")
	}

	verificationCode := generateSixDigitCode()

	payload := map[string]interface{}{
		"from": map[string]string{
			"email": mailtrapFromEmail,
			"name":  mailtrapFromName,
		},
		"to":       []map[string]string{{"email": to}},
		"subject":  "Your two-factor verification code",
		"html":     buildTwoFactorEmailHTML(verificationCode),
		"category": "Two Factor Verification",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, mailtrapAPIURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+mailtrapToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("mailtrap request failed with status: %s", res.Status)
	}

	return verificationCode, nil
}

func generateSixDigitCode() string {
	var n uint32
	_ = binary.Read(rand.Reader, binary.LittleEndian, &n)
	return fmt.Sprintf("%06d", n%1000000)
}

func buildTwoFactorEmailHTML(code string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Your Two-Factor Code</title>
  </head>
  <body style="margin:0;padding:0;background-color:#f4f4f4;font-family:Arial,Helvetica,sans-serif;color:#111111;">
    <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="background-color:#f4f4f4;padding:24px 0;">
      <tr>
        <td align="center">
          <table role="presentation" width="600" cellspacing="0" cellpadding="0" border="0" style="background-color:#ffffff;border:1px solid #e5e5e5;border-radius:8px;overflow:hidden;">
            <tr>
              <td style="padding:32px 32px 16px 32px;">
                <h1 style="margin:0;font-size:24px;color:#111111;">Your Project Name</h1>
                <p style="margin:16px 0 0 0;font-size:16px;line-height:1.6;color:#444444;">Use the code below to complete your login.</p>
              </td>
            </tr>
            <tr>
              <td style="padding:0 32px 16px 32px;">
                <div style="background-color:#f8f8f8;border:1px solid #dddddd;border-radius:6px;padding:16px;text-align:center;font-size:24px;font-weight:bold;letter-spacing:1px;color:#111111;">%s</div>
              </td>
            </tr>
            <tr>
              <td style="padding:0 32px 24px 32px;">
                <p style="margin:0;font-size:14px;line-height:1.6;color:#666666;">This code is valid for a short time. Do not share it with anyone.</p>
              </td>
            </tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>`, code)
}

func buildVerificationEmailHTML(code string) string {
	baseURL := strings.TrimRight(os.Getenv("APP_BASE_URL"), "/")
	if baseURL == "" {
		baseURL = "https://yourdomain.com"
	}

	verificationLink := fmt.Sprintf("%s/verify-account?code=%s", baseURL, url.QueryEscape(code))

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Verify Your Account</title>
  </head>
  <body style="margin:0;padding:0;background-color:#f4f4f4;font-family:Arial,Helvetica,sans-serif;color:#111111;">
    <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="background-color:#f4f4f4;padding:24px 0;">
      <tr>
        <td align="center">
          <table role="presentation" width="600" cellspacing="0" cellpadding="0" border="0" style="background-color:#ffffff;border:1px solid #e5e5e5;border-radius:8px;overflow:hidden;">
            <tr>
              <td style="padding:32px 32px 16px 32px;">
                <h1 style="margin:0;font-size:24px;color:#111111;">Your Project Name</h1>
                <p style="margin:16px 0 0 0;font-size:16px;line-height:1.6;color:#444444;">Please verify your account using the code below.</p>
              </td>
            </tr>
            <tr>
              <td style="padding:0 32px 16px 32px;">
                <div style="background-color:#f8f8f8;border:1px solid #dddddd;border-radius:6px;padding:16px;text-align:center;font-size:24px;font-weight:bold;letter-spacing:1px;color:#111111;">%s</div>
              </td>
            </tr>
            <tr>
              <td style="padding:0 32px 24px 32px;">
                <a href="%s" style="display:inline-block;background-color:#111111;color:#ffffff;text-decoration:none;padding:12px 20px;border-radius:4px;font-size:16px;">Verify Account</a>
              </td>
            </tr>
            <tr>
              <td style="padding:0 32px 32px 32px;">
                <p style="margin:0;font-size:14px;line-height:1.6;color:#666666;">If the button does not work, use this link:</p>
                <p style="margin:8px 0 0 0;font-size:14px;word-break:break-all;"><a href="%s" style="color:#111111;">%s</a></p>
              </td>
            </tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>`, code, verificationLink, verificationLink, verificationLink)
}
