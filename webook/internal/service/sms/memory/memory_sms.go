package memory

import (
	"context"
	"fmt"
	"start/webook/internal/service/sms"
)

type memory struct {
}

func NewMemory() sms.SMS {
	return &memory{}
}

func (m memory) Send(ctx context.Context, phone string, biz string, templateParam []string) error {
	fmt.Println(templateParam[1], "这是m1")
	return nil
}
