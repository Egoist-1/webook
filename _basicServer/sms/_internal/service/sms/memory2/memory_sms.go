package memory

import (
	"context"
	"fmt"
	"webook/_basicServer/sms/_internal/service/sms"
)

type memory2 struct {
}

func NewMemory2() sms.SMS {
	return &memory2{}
}

func (m memory2) Send(ctx context.Context, phone string, biz string, templateParam []string) error {
	fmt.Println(templateParam[1], "这是m2")
	return nil
}
