package sms

import "context"

type SMS interface {
	Send(ctx context.Context, phone string, biz string, templateParam []string) error
}
