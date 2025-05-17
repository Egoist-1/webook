package failover

import (
	"context"
	"webook/_basicServer/sms/_internal/service/sms"
)

type failover struct {
	ss []sms.SMS
}

func NewFailoverSMS(sms ...sms.SMS) sms.SMS {
	return &failover{
		ss: sms,
	}
}

func (a failover) Send(ctx context.Context, phone string, biz string, templateParam []string) error {
	var err error
	for _, sms := range a.ss {
		err = sms.Send(ctx, phone, biz, templateParam)
		if err == nil {
			return nil
		}
	}
	return err
}
