package services

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendInviteEmail(toEmail string, inviteLink string) error
	SendPasswordChangedAlert(toEmail string) error
	SendPasswordResetEmail(toEmail, resetLink string) error
	SendSLAAlert(toEmail, requestID string, waitTimeHours int) error
}

type emailService struct {
	smtpHost string
	smtpPort int
	email    string
	password string
	dialer   *gomail.Dialer
}

func NewEmailService(host, portStr, email, password string) (EmailService, error) {
	host = strings.TrimSpace(host)
	portStr = strings.TrimSpace(portStr)
	email = strings.TrimSpace(email)

	if host == "" {
		return nil, fmt.Errorf("SMTP_HOST is required")
	}
	if portStr == "" {
		return nil, fmt.Errorf("SMTP_PORT is required")
	}
	if email == "" {
		return nil, fmt.Errorf("SMTP_EMAIL is required")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid smtp port: %w", err)
	}
	if port <= 0 || port > 65535 {
		return nil, fmt.Errorf("invalid smtp port: %d", port)
	}

	dialer := gomail.NewDialer(host, port, email, password)

	return &emailService{
		smtpHost: host,
		smtpPort: port,
		email:    email,
		password: password,
		dialer:   dialer,
	}, nil
}

func (s *emailService) SendInviteEmail(toEmail string, inviteLink string) error {
	fromEmail := s.email
	if fromEmail == "" {
		fromEmail = "noreply@Omnilog.local"
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("Omnilog System <%s>", fromEmail))
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Запрошення до системи Omnilog")

	plainText := fmt.Sprintf("Вітаємо!\nВаш акаунт створено.\nПерейдіть за посиланням для активації: %s", inviteLink)
	m.SetBody("text/plain", plainText)

	htmlBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 5px;">
			<h2 style="color: #2c3e50;">Ласкаво просимо до OmniLog</h2>
			<p>Вам надано доступ до системи управління логістикою.</p>
			<p>Для активації облікового запису та встановлення пароля, натисніть на кнопку нижче:</p>
			
			<div style="text-align: center; margin: 30px 0;">
				<a href="%s" style="background-color: #0056b3; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; font-weight: bold;">
					Активувати Акаунт
				</a>
			</div>
			
			<p style="font-size: 12px; color: #7f8c8d;">
				Посилання дійсне протягом 24 годин.<br>
				Якщо кнопка не працює, скопіюйте це посилання у браузер:<br>
				<a href="%s">%s</a>
			</p>
			<hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
			<p style="font-size: 10px; color: #999;">Це автоматичне повідомлення, не відповідайте на нього.</p>
		</div>
	`, inviteLink, inviteLink, inviteLink)

	m.AddAlternative("text/html", htmlBody)

	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email via smtp: %w", err)
	}

	return nil
}

func (s *emailService) SendPasswordChangedAlert(toEmail string) error {
	subject := "Увага: Ваш пароль було змінено"
	body := "Вітаємо!\n\nВаш пароль у системі OmniLog щойно було успішно змінено.\n\nЯкщо ви цього не робили, негайно зверніться до вашого керівника або адміністратора системи."

	return s.sendMail(toEmail, subject, body)
}

// 2. Для Забув пароль (Лінк на відновлення)
func (s *emailService) SendPasswordResetEmail(toEmail, resetLink string) error {
	subject := "Omnilog: Відновлення пароля"
	body := fmt.Sprintf("Ви подали запит на скидання пароля.\n\nПерейдіть за посиланням нижче, щоб встановити новий пароль. Посилання дійсне 1 годину:\n\n%s\n\nЯкщо ви не робили цього запиту, просто проігноруйте цей лист.", resetLink)

	htmlBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 5px;">
			<h2 style="color: #2c3e50;">Відновлення пароля OmniLog</h2>
			<p>Ви подали запит на скидання пароля.</p>
			<p>Щоб встановити новий пароль, натисніть на кнопку нижче:</p>
			<div style="text-align: center; margin: 30px 0;">
				<a href="%s" style="background-color: #0056b3; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; font-weight: bold;">
					Скинути пароль
				</a>
			</div>
			<p style="font-size: 12px; color: #7f8c8d;">
				Посилання дійсне протягом 1 години.<br>
				Якщо кнопка не працює, скопіюйте це посилання у браузер:<br>
				<a href="%s">%s</a>
			</p>
			<hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
			<p style="font-size: 10px; color: #999;">Якщо ви не робили цього запиту, просто проігноруйте цей лист.</p>
		</div>
	`, resetLink, resetLink, resetLink)

	return s.sendMailWithHTML(toEmail, subject, body, htmlBody)
}

// Допоміжний метод для відправки простих текстових листів через gomail
func (s *emailService) sendMail(toEmail, subject, plainText string) error {
	return s.sendMailWithHTML(toEmail, subject, plainText, "")
}

func (s *emailService) sendMailWithHTML(toEmail, subject, plainText, htmlBody string) error {
	fromEmail := s.email
	if fromEmail == "" {
		fromEmail = "noreply@Omnilog.local"
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("Omnilog System <%s>", fromEmail))
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)

	m.SetBody("text/plain", plainText)
	if htmlBody != "" {
		m.AddAlternative("text/html", htmlBody)
	}

	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email via smtp: %w", err)
	}

	return nil
}

func (s *emailService) SendSLAAlert(toEmail, requestID string, waitTimeHours int) error {
	subject := "🚨 Увага: Порушення SLA (Заявка зависла)"
	body := fmt.Sprintf("Автоматичне сповіщення системи OmniLog.\n\nЗаявка #%s очікує підтвердження понад %d годин і порушує стандарти SLA.\n\nБудь ласка, перевірте панель керування та прийміть рішення.", requestID, waitTimeHours)

	return s.sendMail(toEmail, subject, body)
}
