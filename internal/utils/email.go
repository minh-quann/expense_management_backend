package utils

import (
	"fmt"
	"log"
	"net/smtp"

	"expense_management_backend/config"
)

// SendResetPasswordEmail sends a password reset code to the user's email address
func SendResetPasswordEmail(cfg *config.Config, toEmail, resetToken string) error {
	// If SMTP host is not configured, fall back to console logging
	if cfg.SMTP.Host == "" {
		log.Printf("\n🔑 [SMTP Mock] PASSWORD RESET CODE for %s: %s\n", toEmail, resetToken)
		return nil
	}

	// Prepare email headers and HTML body
	subject := "Reset Your Expense Management Password"
	from := cfg.SMTP.From
	
	// HTML template for password reset email
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Reset Password</title>
</head>
<body style="font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; background-color: #f4f6f9; margin: 0; padding: 40px; color: #333333;">
    <div style="max-width: 500px; margin: 0 auto; background-color: #ffffff; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 10px rgba(0, 0, 0, 0.05); border: 1px solid #eef2f5;">
        <div style="background-color: #5A45FE; padding: 30px; text-align: center;">
            <h1 style="color: #ffffff; margin: 0; font-size: 24px; font-weight: 700;">Expense Management</h1>
        </div>
        <div style="padding: 40px 30px;">
            <p style="font-size: 16px; line-height: 1.6; margin-top: 0; margin-bottom: 24px;">Hi,</p>
            <p style="font-size: 16px; line-height: 1.6; margin-bottom: 24px;">We received a request to reset your password. Please use the verification code below to reset it. This code will expire in 15 minutes.</p>
            
            <div style="background-color: #f0edff; border-radius: 8px; padding: 16px; text-align: center; margin-bottom: 24px;">
                <span style="font-size: 32px; font-weight: 700; color: #5A45FE; letter-spacing: 4px;">%s</span>
            </div>
            
            <p style="font-size: 14px; line-height: 1.6; color: #777777; margin-bottom: 0;">If you did not request a password reset, please ignore this email.</p>
        </div>
        <div style="background-color: #fcfdfe; padding: 20px 30px; text-align: center; border-top: 1px solid #eef2f5;">
            <span style="font-size: 12px; color: #999999;">&copy; 2026 Expense Management Inc. All rights reserved.</span>
        </div>
    </div>
</body>
</html>`, resetToken)

	// Format RFC 822 email message
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := []byte("To: " + toEmail + "\n" +
		"From: " + from + "\n" +
		"Subject: " + subject + "\n" +
		mime + body)

	// Auth and address configuration
	addr := fmt.Sprintf("%s:%s", cfg.SMTP.Host, cfg.SMTP.Port)
	var auth smtp.Auth
	if cfg.SMTP.User != "" && cfg.SMTP.Password != "" {
		auth = smtp.PlainAuth("", cfg.SMTP.User, cfg.SMTP.Password, cfg.SMTP.Host)
	}

	// Send email
	err := smtp.SendMail(addr, auth, from, []string{toEmail}, message)
	if err != nil {
		log.Printf("❌ Failed to send reset email to %s: %v\n", toEmail, err)
		return err
	}

	log.Printf("📧 Reset email successfully sent to %s\n", toEmail)
	return nil
}
