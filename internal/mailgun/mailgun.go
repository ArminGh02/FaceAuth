package mailgun

import (
	"log"

	"github.com/mailgun/mailgun-go"
)

type MailGun struct {
	Config Config
	mg     mailgun.Mailgun
}

func New(cfg Config) *MailGun {
	return &MailGun{
		Config: cfg,
		mg:     mailgun.NewMailgun(cfg.Domain, cfg.APIKey),
	}
}

func (m *MailGun) Send(sender, subject, body, recipient string) error {
	message := m.mg.NewMessage(sender, subject, body, recipient)
	resp, id, err := m.mg.Send(message)
	log.Printf("sending email from %s to %s with subject %s and body:\n%s\nresponse: %s id: %s\n", sender, recipient, subject, body, resp, id)
	return err
}
