package email

import (
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridService struct {
	apiKey    string
	fromEmail string
	fromName  string
}

func NewSendGridService(apiKey, fromEmail, fromName string) *SendGridService {
	return &SendGridService{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		fromName:  fromName,
	}
}

func (s *SendGridService) SendBusinessCreated(to, businessName string) error {
	from := mail.NewEmail(s.fromName, s.fromEmail)
	toEmail := mail.NewEmail("", to)
	subject := "Clínica creada"
	content := mail.NewContent("text/html", fmt.Sprintf("<h1>Tu clínica %s fue creada</h1>", businessName))
	message := mail.NewV3MailInit(from, subject, toEmail, content)

	client := sendgrid.NewSendClient(s.apiKey)
	resp, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("sendgrid send: %w", err)
	}

	log.Printf("sendgrid response: status=%d body=%s", resp.StatusCode, resp.Body)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("sendgrid rejected email: status=%d body=%s", resp.StatusCode, resp.Body)
	}

	return nil
}
