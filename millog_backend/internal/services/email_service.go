package services

import (
	"fmt"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendInviteEmail(toEmail string, inviteLink string) error
}

type emailService struct {
	smtpHost string
	smtpPort int
	email    string
	password string
	dialer   *gomail.Dialer
}

func NewEmailService(host, portStr, email, password string) (EmailService, error) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid smtp port: %w", err)
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
	m := gomail.NewMessage()

	m.SetHeader("From", fmt.Sprintf("Millog System <%s>", s.email))
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Запрошення до системи Millog")

	plainText := fmt.Sprintf("Вітаємо!\nВаш акаунт створено.\nПерейдіть за посиланням для активації: %s", inviteLink)
	m.SetBody("text/plain", plainText)

	htmlBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 5px;">
			<h2 style="color: #2c3e50;">Ласкаво просимо до Millog</h2>
			<p>Вам надано доступ до системи військової логістики.</p>
			<p>Для активації облікового запису та встановлення пароля, натисніть на кнопку нижче:</p>
			
			<div style="text-align: center; margin: 30px 0;">
				<a href="%s" style="background-color: #0056b3; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; font-weight: bold;">
					Активувати Акаунт
				</a>
			</div>
			
			<p style="font-size: 12px; color: #7f8c8d;">
				Посилання дійсне протягом 48 годин.<br>
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
