package article

import (
	"context"
)

type ArticleDao interface {
	UpdateById(ctx context.Context, art Article) (int, error)
	Create(ctx context.Context, art Article) (int, error)
	Sync(ctx context.Context, entity Article) (int, error)
	GetUnpublishList(ctx context.Context, id int, limit int, offset int) (arts []Article, err error)
	FindById(ctx context.Context, uid int, aid int) (Article, error)
	PubFinById(ctx context.Context, uid int, aid int, status uint) (Article, error)
}
