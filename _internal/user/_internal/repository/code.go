package repository

import (
	"context"
	"webook/_internal/user/_internal/repository/cache"
)

type CodeRepo interface {
	Store(ctx context.Context, key string, val string) error
	Verify(ctx context.Context, key string, code string) error
}

type codeRepo struct {
	cache cache.CodeCache
}

func NewCodeRepo(cache cache.CodeCache) CodeRepo {
	return &codeRepo{cache: cache}
}

func (c *codeRepo) Verify(ctx context.Context, key string, code string) error {
	return c.cache.Verify(ctx, key, code)
}

func (c *codeRepo) Store(ctx context.Context, key string, val string) error {
	return c.cache.Store(ctx, key, val)
}
