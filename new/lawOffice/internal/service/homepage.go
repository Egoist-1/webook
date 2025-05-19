package service

import (
	"errors"
	"lawOffice/internal/e"
)

func A() error {
	err := e.NewErr(e.UserExist, "这里错误123", errors.New("").Error())
	return err
}
