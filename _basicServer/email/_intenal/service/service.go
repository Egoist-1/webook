package service

import (
	"context"
	_ "embed"
	"webook/pkg/emailx"
)

type EmailService interface {
	Send(ctx context.Context, email string, biz string, params []string) error
}

var (
	//go:embed template/verifyCode.html
	code string
)

type emailService struct {
	email emailx.EmailClient
}

func (e *emailService) Send(ctx context.Context, email string, biz string, params []string) error {
	err := e.email.Send(ctx, emailx.Req{
		ToEmail: []string{email},
		Subject: "验证码",
		Text:    "",
		HTML:    code,
	})
	return err
}
