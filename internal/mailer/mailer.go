package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
)

//go:embed "templates"
var templateFS embed.FS

// Mailer encapsulates mail config options.
type Mailer struct {
	dialer *mail.Dialer
	sender string
}

// New returns a configured mailer instance.
func New(port int, host, sender, username, password string) Mailer {
	dailer := mail.NewDialer(host, port, username, password)
	dailer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dailer,
		sender: sender,
	}
}

// Send is used for configuring an email template and sending the email to the targeted recipient.
func (m *Mailer) Send(data interface{}, templateFile, reciepient string) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", reciepient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	for i := 0; i < 3; i++ {
		err = m.dialer.DialAndSend(msg)
		if err == nil {
			return err
		}

		time.Sleep(500 * time.Millisecond)
	}

	return err
}
