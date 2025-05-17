package cache

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"webook/pkg/er"
)

type CodeCache interface {
	Store(ctx context.Context, key, val string) error
	Verify(ctx context.Context, key, val string) error
}

type codeCacheRedis struct {
	rdb redis.Cmdable
}

func NewCodeCacheRedis(rdb redis.Cmdable) CodeCache {
	return &codeCacheRedis{rdb: rdb}
}

//go:embed lua/store_verification_code.lua
var store_lua string

func (c codeCacheRedis) Store(ctx context.Context, key, val string) error {
	i, err := c.rdb.Eval(ctx, store_lua, []string{key}, val).Int()
	if err != nil {
		return err
	}
	switch i {
	case 0:
		return nil
	case -1:
		return er.NewErr(er.ServerErr, "codeCacheRedis 没有设置过期时间", "")
	case -2:
		return er.NewErr(er.UserOperationTooFrequent, "codeCacheRedis 验证码重复发送", "")
	}
	return err
}

//go:embed lua/verify_code.lua
var verify_lua string

func (c codeCacheRedis) Verify(ctx context.Context, key, val string) error {
	i, err := c.rdb.Eval(ctx, verify_lua, []string{key}, val).Int()
	if err != nil {
		return err
	}
	switch i {
	case 0:
		return nil
	case -1:
		//验证次数过多
		return er.NewErr(er.Code_TooManyVerificationAttempts, "验证次数过多", "")
	case -2:
		//验证码错误
		return er.NewErr(er.Code_VerifyFail, "验证码错误", "")
	case -3:
		//KEY不存在
		return er.NewErr(er.Code_NotFind, "Key不存在", "")
	}
	return err
}
