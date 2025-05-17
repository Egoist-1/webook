package memory

import (
	"context"
	"fmt"
	"webook/_basicServer/sms/_internal/service/sms"
)

type memory3 struct {
}

func NewMemory3() sms.SMS {
	return &memory3{}
}

func (m memory3) Send(ctx context.Context, phone string, biz string, templateParam []string) error {
	fmt.Println(templateParam[1], "这是m3")
	return nil
}
