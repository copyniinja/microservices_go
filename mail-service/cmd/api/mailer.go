package main

import (
	"bytes"
	"text/template"
	"time"

	gomail "github.com/go-mail/mail/v2"
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}
type Mailer struct {
	Mail
	Dialer *gomail.Dialer
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

func NewMailer(m Mail) *Mailer {
	d := gomail.NewDialer(m.Host, m.Port, m.Username, m.Password)

	d.Timeout = 10 * time.Second

	switch m.Encryption {
	case "ssl":
		d.SSL = true
	case "tls":
		d.StartTLSPolicy = gomail.MandatoryStartTLS
	default:
		d.StartTLSPolicy = gomail.NoStartTLS
	}
	return &Mailer{
		Mail:   m,
		Dialer: d,
	}
}

func (m *Mailer) Send(templateFile string, msg Message) error {

	// Parse template
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	var body bytes.Buffer

	data := msg.Data
	if msg.DataMap != nil {
		data = msg.DataMap
	}

	if err := t.Execute(&body, data); err != nil {
		return err
	}

	// Build email
	email := gomail.NewMessage()

	from := m.FromAddress
	fromName := m.FromName

	if msg.From != "" {
		from = msg.From
	}

	if msg.FromName != "" {
		fromName = msg.FromName
	}

	email.SetHeader("From", email.FormatAddress(from, fromName))
	email.SetHeader("To", msg.To)
	email.SetHeader("Subject", msg.Subject)

	email.SetBody("text/html", body.String())

	// Attachments
	for _, file := range msg.Attachments {
		email.Attach(file)
	}

	// Send
	return m.Dialer.DialAndSend(email)
}
