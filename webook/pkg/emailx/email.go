package emailx

import (
	"context"
	"github.com/jordan-wright/email"
)

type EmailClient interface {
	Send(ctx context.Context, req Req) error
}
type Req struct {
	ToEmail []string
	Subject string
	Text    string
	HTML    string
}

func NewEmail(fromEmail string) EmailClient {
	e := email.NewEmail()
	e.From = fromEmail
	return &emailx{
		FromEmail: fromEmail,
		client:    e,
	}
}

type emailx struct {
	FromEmail string
	client    *email.Email
}

func (e *emailx) Send(ctx context.Context, req Req) error {
	e.client.To = req.ToEmail
	e.client.Subject = req.Subject
	e.client.Text = []byte(req.Text)
	e.client.HTML = []byte(req.HTML)
	err := e.client.Send("", nil)
	return err
}
