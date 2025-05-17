package service

import (
	"context"
	"go.uber.org/zap"
	"webook/_internal/article/_internal/domain"
	"webook/_internal/article/_internal/event/article"
	"webook/_internal/article/_internal/repository"
)

//go:generate mockgen -source=article.go -package=svcmocks -destination=mocks/article.mock.go ArticleService
type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	List(ctx context.Context, uid int64, limit int, offset int) ([]domain.Article, error)
	Detail(ctx context.Context, uid int64, aid int64) (domain.Article, error)
	PubDetail(ctx context.Context, id int64, aid int64) (domain.Article, error)
	TopN(ctx context.Context, ids []int64) ([]domain.Article, error)
	PubList(ctx context.Context, uid int64, limit int, offset int) ([]domain.Article, error)
	Unpublish(ctx context.Context, aid int64) error
}

func NewArticleService(repo repository.ArticleRepository, producer article.Producer) ArticleService {
	return &articleService{
		producer: producer,
		repo:     repo,
	}
}

type articleService struct {
	repo     repository.ArticleRepository
	producer article.Producer
}

func (svc *articleService) Unpublish(ctx context.Context, aid int64) error {
	return svc.repo.SyncStatus(ctx, aid)
}

func (svc *articleService) PubList(ctx context.Context, uid int64, limit int, offset int) ([]domain.Article, error) {
	return svc.repo.PubList(ctx, uid, limit, offset)
}

func (svc *articleService) TopN(ctx context.Context, ids []int64) ([]domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (svc *articleService) Detail(ctx context.Context, uid int64, aid int64) (domain.Article, error) {
	return svc.repo.Detail(ctx, uid, aid)
}
func (svc *articleService) PubDetail(ctx context.Context, uid int64, aid int64) (domain.Article, error) {
	art, err := svc.repo.PubDetail(ctx, uid, aid, domain.ArticleStatusPublished)
	if err != nil {
		return domain.Article{}, err
	}

	go func() {
		er := svc.producer.IncrReadCnt(ctx, article.ReadEvent{
			Uid: art.Id,
			Aid: art.Author.Id,
		})
		if er != nil {
			zap.L().Error("消息队列发送失败", zap.Error(er))
		}
	}()
	return art, err
}

func (svc *articleService) List(ctx context.Context, id int64, limit int, offset int) ([]domain.Article, error) {
	return svc.repo.List(ctx, id, limit, offset)
}

func (svc *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return svc.repo.Sync(ctx, art)
}

func (svc *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnpublished
	if art.Id > 0 {
		return svc.repo.Update(ctx, art)
	}
	id, err := svc.repo.Create(ctx, art)
	return id, err
}
