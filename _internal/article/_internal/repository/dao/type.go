package dao

import (
	"context"
)

type ArticleDao interface {
	UpdateById(ctx context.Context, art Article) (int64, error)
	Create(ctx context.Context, art Article) (int64, error)
	Sync(ctx context.Context, entity Article) (int64, error)
	GetUnpublishList(ctx context.Context, id int64, limit int, offset int) (arts []Article, err error)
	FindById(ctx context.Context, uid int64, aid int64) (Article, error)
	PubFinById(ctx context.Context, uid int64, aid int64, status uint) (Article, error)
	PubList(ctx context.Context, uid int64, limit int, offset int) ([]ArticlePublish, error)
	SyncStatus(ctx context.Context, aid int64) error
}
