package retry

import (
	"context"
	"webook/_basicServer/sms/_internal/service/sms"
)

type retry struct {
	sms.SMS
	retryMax int
}

func NewRetry(SMS sms.SMS) sms.SMS {
	return &retry{
		SMS:      SMS,
		retryMax: 3,
	}
}
func (m retry) Send(ctx context.Context, phone string, biz string, templateParam []string) error {
	err := m.Send(ctx, phone, biz, templateParam)
	if err == nil {
		return nil
	}
	for _ = range m.retryMax {
		err = m.Send(ctx, phone, biz, templateParam)
		if err != nil {
			continue
		}
	}
	return err
}
