package service

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/bayupaths/bypur-api/internal/config"
)

type MailService struct {
	cfg *config.Config
}

func NewMailService(cfg *config.Config) *MailService {
	return &MailService{cfg: cfg}
}

func (s *MailService) VerifyConnection() error {
	if !s.cfg.Mail.Enabled {
		return fmt.Errorf("SMTP is disabled")
	}

	addr := s.cfg.Mail.Host + ":" + strconv.Itoa(s.cfg.Mail.Port)
	isSSL := s.cfg.Mail.Port == 465

	var conn net.Conn
	var err error

	if isSSL {
		conn, err = tls.Dial("tcp", addr, &tls.Config{
			ServerName: s.cfg.Mail.Host,
		})
	} else {
		dialer := net.Dialer{Timeout: 5 * time.Second}
		conn, err = dialer.Dial("tcp", addr)
	}
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.cfg.Mail.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	if err = client.Hello("localhost"); err != nil {
		return fmt.Errorf("hello command failed: %w", err)
	}

	if !isSSL {
		if ok, _ := client.Extension("STARTTLS"); ok {
			tlsConfig := &tls.Config{ServerName: s.cfg.Mail.Host}
			if err = client.StartTLS(tlsConfig); err != nil {
				return fmt.Errorf("STARTTLS handshake failed: %w", err)
			}
		}
	}

	if s.cfg.Mail.User != "" && s.cfg.Mail.Pass != "" {
		auth := smtp.PlainAuth("", s.cfg.Mail.User, s.cfg.Mail.Pass, s.cfg.Mail.Host)
		if ok, _ := client.Extension("AUTH"); ok {
			if err = client.Auth(auth); err != nil {
				return fmt.Errorf("SMTP authentication failed: %w", err)
			}
		}
	}

	return nil
}

func (s *MailService) SendEmail(to, subject, body string, replyTo ...string) error {
	if !s.cfg.Mail.Enabled {
		slog.Warn("Attempted to send email but SMTP is disabled", "to", to)
		return nil
	}

	addr := s.cfg.Mail.Host + ":" + strconv.Itoa(s.cfg.Mail.Port)
	isSSL := s.cfg.Mail.Port == 465

	var conn net.Conn
	var err error

	if isSSL {
		conn, err = tls.Dial("tcp", addr, &tls.Config{
			ServerName: s.cfg.Mail.Host,
		})
	} else {
		dialer := net.Dialer{Timeout: 10 * time.Second}
		conn, err = dialer.Dial("tcp", addr)
	}
	if err != nil {
		slog.Error("Failed to connect to SMTP server", "addr", addr, "error", err)
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.cfg.Mail.Host)
	if err != nil {
		slog.Error("Failed to create SMTP client", "error", err)
		return err
	}
	defer client.Close()

	if err = client.Hello("localhost"); err != nil {
		slog.Error("SMTP Hello failed", "error", err)
		return err
	}

	if !isSSL {
		if ok, _ := client.Extension("STARTTLS"); ok {
			tlsConfig := &tls.Config{ServerName: s.cfg.Mail.Host}
			if err = client.StartTLS(tlsConfig); err != nil {
				slog.Error("SMTP STARTTLS failed", "error", err)
				return err
			}
		}
	}

	if s.cfg.Mail.User != "" && s.cfg.Mail.Pass != "" {
		auth := smtp.PlainAuth("", s.cfg.Mail.User, s.cfg.Mail.Pass, s.cfg.Mail.Host)
		if ok, _ := client.Extension("AUTH"); ok {
			if err = client.Auth(auth); err != nil {
				slog.Error("SMTP Auth failed", "error", err)
				return err
			}
		}
	}

	if err = client.Mail(s.cfg.Mail.From); err != nil {
		slog.Error("SMTP Mail command failed", "from", s.cfg.Mail.From, "error", err)
		return err
	}

	if err = client.Rcpt(to); err != nil {
		slog.Error("SMTP Rcpt command failed", "to", to, "error", err)
		return err
	}

	w, err := client.Data()
	if err != nil {
		slog.Error("SMTP Data command failed", "error", err)
		return err
	}

	headerTo := to
	headerFrom := s.cfg.Mail.From

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

	_, err = w.Write([]byte(message))
	if err != nil {
		slog.Error("Failed to write email body", "error", err)
		w.Close()
		return err
	}

	err = w.Close()
	if err != nil {
		slog.Error("Failed to close data writer", "error", err)
		return err
	}

	err = client.Quit()
	if err != nil {
		slog.Warn("SMTP Quit command returned error", "error", err)
	}

	slog.Info("Email sent successfully", "to", to, "subject", subject)
	return nil
}

func (s *MailService) SendVerificationEmail(to, token string) error {
	verificationUrl := s.cfg.Frontend.CmsURL + "/verify-email?token=" + token
	body := s.getVerificationEmailTemplate(verificationUrl)
	return s.SendEmail(to, "Verify Your Email Address", body)
}

func (s *MailService) SendPasswordResetEmail(to, token string) error {
	resetUrl := s.cfg.Frontend.CmsURL + "/reset-password?token=" + token
	body := s.getPasswordResetEmailTemplate(resetUrl)
	return s.SendEmail(to, "Reset Your Password", body)
}

func (s *MailService) SendWelcomeEmail(to, name string) error {
	body := s.getWelcomeEmailTemplate(name)
	return s.SendEmail(to, "Welcome to "+s.cfg.App.Name+"!", body)
}

func (s *MailService) SendContactFormEmail(name, email, message string) error {
	body := s.getContactFormEmailTemplate(name, email, message)
	return s.SendEmail(s.cfg.Mail.From, "New Contact Form Submission from "+name, body, email)
}

func (s *MailService) getHeaderLogo() (string, string) {
	name := s.cfg.App.Name
	if name == "" {
		name = "bayu-apps"
	}
	if strings.Contains(name, "-") {
		parts := strings.SplitN(name, "-", 2)
		return parts[0], "-" + parts[1]
	}
	if strings.Contains(name, " ") {
		parts := strings.SplitN(name, " ", 2)
		return parts[0], " " + parts[1]
	}
	return name, ""
}

func (s *MailService) wrapEmailContent(innerContent string) string {
	year := time.Now().Year()
	logoPrefix, logoSuffix := s.getHeaderLogo()

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <style>
      body {
        background-color: #f8fafc;
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
        -webkit-font-smoothing: antialiased;
        font-size: 14px;
        line-height: 1.6;
        margin: 0;
        padding: 0;
        -ms-text-size-adjust: 100%%;
        -webkit-text-size-adjust: 100%%;
      }
      .wrapper {
        background-color: #f8fafc;
        width: 100%%;
        padding: 40px 0;
      }
      .container {
        background-color: #ffffff;
        border: 1px solid #e2e8f0;
        border-radius: 8px;
        box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.025);
        max-width: 580px;
        margin: 0 auto;
        padding: 40px;
      }
      .header {
        margin-bottom: 30px;
        text-align: center;
      }
      .header h1 {
        color: #0f172a;
        font-size: 24px;
        font-weight: 700;
        margin: 0;
        letter-spacing: -0.025em;
      }
      .header-accent {
        color: #0ea5e9;
      }
      .content {
        color: #334155;
      }
      .content h2 {
        color: #0f172a;
        font-size: 20px;
        font-weight: 600;
        margin-top: 0;
        margin-bottom: 16px;
      }
      .content p {
        margin-top: 0;
        margin-bottom: 16px;
        color: #475569;
        font-size: 15px;
      }
      .button-container {
        margin: 30px 0;
        text-align: center;
      }
      .btn {
        background-color: #032448;
        border-radius: 6px;
        color: #ffffff !important;
        display: inline-block;
        font-size: 15px;
        font-weight: 600;
        line-height: 2.4;
        padding: 4px 24px;
        text-align: center;
        text-decoration: none;
        box-shadow: 0 2px 4px rgba(3, 36, 72, 0.2);
      }
      .btn-secondary {
        background-color: #0ea5e9;
        box-shadow: 0 2px 4px rgba(14, 165, 233, 0.2);
      }
      .divider {
        border-top: 1px solid #e2e8f0;
        margin: 24px 0;
      }
      .link-fallback {
        background-color: #f1f5f9;
        border-radius: 6px;
        font-family: monospace;
        font-size: 12px;
        padding: 12px;
        word-break: break-all;
        color: #64748b;
      }
      .footer {
        color: #94a3b8;
        font-size: 12px;
        text-align: center;
        margin-top: 30px;
      }
      .footer p {
        margin: 4px 0;
      }
      .footer a {
        color: #0ea5e9;
        text-decoration: none;
      }
      .form-table {
        width: 100%%;
        border-collapse: collapse;
        margin: 20px 0;
      }
      .form-table td {
        padding: 12px;
        border-bottom: 1px solid #f1f5f9;
      }
      .form-table td.label {
        font-weight: 600;
        color: #1e293b;
        width: 100px;
        vertical-align: top;
      }
      .form-table td.value {
        color: #475569;
      }
      .message-box {
        background-color: #f8fafc;
        border-left: 4px solid #0ea5e9;
        padding: 16px;
        margin: 10px 0;
        font-style: italic;
        border-radius: 0 6px 6px 0;
      }
    </style>
  </head>
  <body>
    <div class="wrapper">
      <div class="container">
        <div class="header">
          <h1>%s<span class="header-accent">%s</span></h1>
        </div>
        <div class="content">
          %s
        </div>
        <div class="footer">
          <p>&copy; %d %s. All rights reserved.</p>
          <p>This is an automated message, please do not reply directly to this email.</p>
        </div>
      </div>
    </div>
  </body>
</html>`, logoPrefix, logoSuffix, innerContent, year, s.cfg.App.Name)
}

func (s *MailService) getVerificationEmailTemplate(verificationUrl string) string {
	inner := fmt.Sprintf(`
          <h2>Email Verification</h2>
          <p>Thank you for signing up! Please verify your email address by clicking the button below:</p>
          <div class="button-container">
            <a href="%s" class="btn" style="color: #ffffff;">Verify Email</a>
          </div>
          <p>Or copy and paste this link in your browser:</p>
          <div class="link-fallback">%s</div>
          <div class="divider"></div>
          <p style="font-size: 13px; color: #94a3b8;">This link will expire in 24 hours. If you didn't create this account, please ignore this email.</p>
	`, verificationUrl, verificationUrl)
	return s.wrapEmailContent(inner)
}

func (s *MailService) getPasswordResetEmailTemplate(resetUrl string) string {
	inner := fmt.Sprintf(`
          <h2>Reset Your Password</h2>
          <p>We received a request to reset your password. Click the button below to set a new password:</p>
          <div class="button-container">
            <a href="%s" class="btn btn-secondary" style="color: #ffffff;">Reset Password</a>
          </div>
          <p>Or copy and paste this link in your browser:</p>
          <div class="link-fallback">%s</div>
          <div class="divider"></div>
          <p style="font-size: 13px; color: #94a3b8;">This link will expire in 1 hour. If you didn't request a password reset, please ignore this email.</p>
	`, resetUrl, resetUrl)
	return s.wrapEmailContent(inner)
}

func (s *MailService) getWelcomeEmailTemplate(name string) string {
	inner := fmt.Sprintf(`
          <h2>Welcome, %s!</h2>
          <p>Thank you for joining us. We're excited to have you on board.</p>
          <p>Our platform offers premium portfolio showcases and admin management features. We hope you enjoy the experience!</p>
          <p>If you have any questions, feel free to reach out to our team.</p>
	`, name)
	return s.wrapEmailContent(inner)
}

func (s *MailService) getContactFormEmailTemplate(name, email, message string) string {
	msgHTML := strings.ReplaceAll(message, "\n", "<br>")
	inner := fmt.Sprintf(`
          <h2>New Contact Form Submission</h2>
          <p>You have received a new message from your portfolio contact form:</p>
          <table class="form-table">
            <tr>
              <td class="label">Name</td>
              <td class="value">%s</td>
            </tr>
            <tr>
              <td class="label">Email</td>
              <td class="value"><a href="mailto:%s" style="color: #0ea5e9; text-decoration: none;">%s</a></td>
            </tr>
          </table>
          <div class="label" style="margin-top: 15px; font-weight: 600; color: #1e293b;">Message:</div>
          <div class="message-box">
            %s
          </div>
          <div class="divider"></div>
          <p style="font-size: 13px; color: #94a3b8;">This message was submitted via the contact form on your portfolio website.</p>
	`, name, email, email, msgHTML)
	return s.wrapEmailContent(inner)
}
