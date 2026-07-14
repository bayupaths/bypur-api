package service

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"bayupur-portofolio-be/internal/config"
)

type MailService struct {
	cfg *config.Config
}

func NewMailService(cfg *config.Config) *MailService {
	return &MailService{cfg: cfg}
}

func (s *MailService) VerifyConnection() error {
	if !s.cfg.SMTPEnabled {
		return fmt.Errorf("SMTP is disabled")
	}

	addr := s.cfg.SMTPHost + ":" + strconv.Itoa(s.cfg.SMTPPort)
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)

	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	if err = client.Hello("localhost"); err != nil {
		return err
	}

	if ok, _ := client.Extension("STARTTLS"); ok {
		slog.Debug("STARTTLS is supported by the SMTP server")
	}

	if err = client.Auth(auth); err != nil {
		return err
	}

	return nil
}

func (s *MailService) SendEmail(to, subject, body string, replyTo ...string) error {
	if !s.cfg.SMTPEnabled {
		slog.Warn("Attempted to send email but SMTP is disabled", "to", to)
		return nil
	}

	addr := s.cfg.SMTPHost + ":" + strconv.Itoa(s.cfg.SMTPPort)
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)

	headerTo := to
	headerFrom := s.cfg.SMTPFrom

	var headers []string
	headers = append(headers, fmt.Sprintf("To: %s", headerTo))
	headers = append(headers, fmt.Sprintf("From: %s", headerFrom))
	headers = append(headers, fmt.Sprintf("Subject: %s", subject))
	headers = append(headers, "MIME-Version: 1.0")
	headers = append(headers, "Content-Type: text/html; charset=UTF-8")

	if len(replyTo) > 0 {
		headers = append(headers, fmt.Sprintf("Reply-To: %s", replyTo[0]))
	}

	message := strings.Join(headers, "\r\n") + "\r\n\r\n" + body

	err := smtp.SendMail(addr, auth, s.cfg.SMTPFrom, []string{to}, []byte(message))
	if err != nil {
		slog.Error("Failed to send email", "to", to, "error", err)
		return err
	}

	slog.Info("Email sent successfully", "to", to, "subject", subject)
	return nil
}

func (s *MailService) SendVerificationEmail(to, token string) error {
	verificationUrl := s.cfg.CmsURL + "/verify-email?token=" + token
	body := s.getVerificationEmailTemplate(verificationUrl)
	return s.SendEmail(to, "Verify Your Email Address", body)
}

func (s *MailService) SendPasswordResetEmail(to, token string) error {
	resetUrl := s.cfg.CmsURL + "/reset-password?token=" + token
	body := s.getPasswordResetEmailTemplate(resetUrl)
	return s.SendEmail(to, "Reset Your Password", body)
}

func (s *MailService) SendWelcomeEmail(to, name string) error {
	body := s.getWelcomeEmailTemplate(name)
	return s.SendEmail(to, "Welcome to "+s.cfg.AppName+"!", body)
}

func (s *MailService) SendContactFormEmail(name, email, message string) error {
	body := s.getContactFormEmailTemplate(name, email, message)
	return s.SendEmail(s.cfg.SMTPFrom, "New Contact Form Submission from "+name, body, email)
}

// Templates HTML
func (s *MailService) getVerificationEmailTemplate(verificationUrl string) string {
	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
      <head>
        <style>
          body { font-family: Arial, sans-serif; }
          .container { max-width: 600px; margin: 0 auto; padding: 20px; }
          .button {
            display: inline-block;
            padding: 12px 24px;
            background-color: #007bff;
            color: white;
            text-decoration: none;
            border-radius: 4px;
          }
          .footer { margin-top: 20px; color: #666; font-size: 12px; }
        </style>
      </head>
      <body>
        <div class="container">
          <h2>Email Verification</h2>
          <p>Thank you for signing up! Please verify your email address by clicking the link below:</p>
          <p>
            <a href="%s" class="button" style="color: white;">Verify Email</a>
          </p>
          <p>Or copy and paste this link in your browser:</p>
          <p>%s</p>
          <div class="footer">
            <p>This link will expire in 24 hours.</p>
            <p>If you didn't create this account, please ignore this email.</p>
          </div>
        </div>
      </body>
    </html>
  `, verificationUrl, verificationUrl)
}

func (s *MailService) getPasswordResetEmailTemplate(resetUrl string) string {
	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
      <head>
        <style>
          body { font-family: Arial, sans-serif; }
          .container { max-width: 600px; margin: 0 auto; padding: 20px; }
          .button {
            display: inline-block;
            padding: 12px 24px;
            background-color: #28a745;
            color: white;
            text-decoration: none;
            border-radius: 4px;
          }
          .footer { margin-top: 20px; color: #666; font-size: 12px; }
        </style>
      </head>
      <body>
        <div class="container">
          <h2>Reset Your Password</h2>
          <p>We received a request to reset your password. Click the link below to create a new password:</p>
          <p>
            <a href="%s" class="button" style="color: white;">Reset Password</a>
          </p>
          <p>Or copy and paste this link in your browser:</p>
          <p>%s</p>
          <div class="footer">
            <p>This link will expire in 1 hour.</p>
            <p>If you didn't request a password reset, please ignore this email.</p>
          </div>
        </div>
      </body>
    </html>
  `, resetUrl, resetUrl)
}

func (s *MailService) getWelcomeEmailTemplate(name string) string {
	year := time.Now().Year()
	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
      <head>
        <style>
          body { font-family: Arial, sans-serif; }
          .container { max-width: 600px; margin: 0 auto; padding: 20px; }
          .footer { margin-top: 20px; color: #666; font-size: 12px; }
        </style>
      </head>
      <body>
        <div class="container">
          <h2>Welcome, %s!</h2>
          <p>Thank you for joining us. We're excited to have you on board.</p>
          <p>If you have any questions, feel free to reach out to our support team.</p>
          <div class="footer">
            <p>&copy; %d %s. All rights reserved.</p>
          </div>
        </div>
      </body>
    </html>
  `, name, year, s.cfg.AppName)
}

func (s *MailService) getContactFormEmailTemplate(name, email, message string) string {
	msgHTML := strings.ReplaceAll(message, "\n", "<br>")
	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
      <head>
        <style>
          body { font-family: Arial, sans-serif; }
          .container { max-width: 600px; margin: 0 auto; padding: 20px; }
          .field { margin: 10px 0; }
          .label { font-weight: bold; }
          .footer { margin-top: 20px; color: #666; font-size: 12px; }
        </style>
      </head>
      <body>
        <div class="container">
          <h2>New Contact Form Submission</h2>
          <div class="field">
            <div class="label">Name:</div>
            <div>%s</div>
          </div>
          <div class="field">
            <div class="label">Email:</div>
            <div>%s</div>
          </div>
          <div class="field">
            <div class="label">Message:</div>
            <div>%s</div>
          </div>
          <div class="footer">
            <p>This message was sent from the contact form.</p>
          </div>
        </div>
      </body>
    </html>
  `, name, email, msgHTML)
}
