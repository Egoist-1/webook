package service

import (
	"context"
	"go.uber.org/zap"
	"start/webook/article/_internal/repository"
	"start/webook/article/internal/domain"
	"start/webook/article/internal/event/article"
)

//go:generate mockgen -source=article.go -package=svcmocks -destination=mocks/article.mock.go ArticleService
type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int, error)
	Publish(ctx context.Context, art domain.Article) (int, error)
	List(ctx context.Context, id int, limit int, offset int) ([]domain.Article, error)
	Detail(ctx context.Context, uid int, aid int) (domain.Article, error)
	PubDetail(ctx context.Context, id int, aid int) (domain.Article, error)
	TopN(ctx context.Context, ids []int64) ([]domain.Article, error)
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

func (svc *articleService) TopN(ctx context.Context, ids []int64) ([]domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (svc *articleService) Detail(ctx context.Context, uid int, aid int) (domain.Article, error) {
	return svc.repo.Detail(ctx, uid, aid)
}
func (svc *articleService) PubDetail(ctx context.Context, uid int, aid int) (domain.Article, error) {
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

func (svc *articleService) List(ctx context.Context, id int, limit int, offset int) ([]domain.Article, error) {
	return svc.repo.List(ctx, id, limit, offset)
}

func (svc *articleService) Publish(ctx context.Context, art domain.Article) (int, error) {
	art.Status = domain.ArticleStatusPublished
	return svc.repo.Sync(ctx, art)
}

func (svc *articleService) Save(ctx context.Context, art domain.Article) (int, error) {
	art.Status = domain.ArticleStatusUnpublished
	if art.Id > 0 {
		return svc.repo.Update(ctx, art)
	}
	id, err := svc.repo.Create(ctx, art)
	return id, err
}
