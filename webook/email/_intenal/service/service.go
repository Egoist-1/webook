package service

import (
	"context"
	"start/webook/pkg/emailx"
)

type EmailService interface {
	Send(ctx context.Context, email string, biz string, params []string) error
}
type emailService struct {
	email emailx.EmailClient
}

func (e *emailService) Send(ctx context.Context, email string, biz string, params []string) error {
	err := e.email.Send(ctx, emailx.Req{
		ToEmail: []string{email},
		Subject: "",
		Text:    "",
		HTML:    "",
	})
	return err
}
