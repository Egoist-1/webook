package service

import (
	"context"
	"math/rand"
	"start/webook/internal/repository"
	"start/webook/internal/service/sms"
)

type CodeService interface {
	SendSMS(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz string, phone string, code string) error
}

type codeService struct {
	s    sms.SMS
	repo repository.CodeRepo
}

func (c codeService) Verify(ctx context.Context, biz string, phone string, code string) error {
	return c.repo.Verify(ctx, c.generateKey(phone, biz), code)
}

func (c codeService) SendSMS(ctx context.Context, phone, biz string) (err error) {
	params := make([]string, 1)
	params[1] = string(rand.Intn(900000) + 100000)
	err = c.s.Send(ctx, phone, biz, params)
	if err != nil {
		return err
	}
	err = c.repo.Store(ctx, c.generateKey(phone, biz), params[1])
	return err
}

func NewCodeService(sms sms.SMS, repo repository.CodeRepo) CodeService {
	return &codeService{
		s:    sms,
		repo: repo,
	}
}

func (c codeService) generateKey(phone, biz string) string {
	return biz + phone
}
