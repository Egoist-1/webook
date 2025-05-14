package repository

import (
	"context"
	"start/webook/article/_internal/repository/dao"
	"start/webook/article/internal/domain"
	"start/webook/pkg/tools/slicex"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int, error)
	Update(ctx context.Context, art domain.Article) (int, error)
	Sync(ctx context.Context, art domain.Article) (int, error)
	List(ctx context.Context, id int, limit int, offset int) ([]domain.Article, error)
	Detail(ctx context.Context, uid int, aid int) (domain.Article, error)
	PubDetail(ctx context.Context, uid int, aid int, status domain.ArticleStatus) (domain.Article, error)
}

func NewArticleRepository(dao dao.ArticleDao) ArticleRepository {
	return &articleRepository{dao: dao}
}

type articleRepository struct {
	dao dao.ArticleDao
}

func (repo *articleRepository) PubDetail(ctx context.Context, uid int, aid int, status domain.ArticleStatus) (domain.Article, error) {
	art, err := repo.dao.PubFinById(ctx, uid, aid, uint(status))
	domainArt := repo.toDomain(art)
	return domainArt, err
}

func (repo *articleRepository) Detail(ctx context.Context, uid int, aid int) (domain.Article, error) {
	art, err := repo.dao.FindById(ctx, uid, aid)
	domainArt := repo.toDomain(art)
	return domainArt, err
}

func (repo *articleRepository) List(ctx context.Context, id int, limit int, offset int) ([]domain.Article, error) {
	arts, err := repo.dao.GetUnpublishList(ctx, id, limit, offset)
	if err != nil {
		return nil, err
	}
	rs := slicex.SliceMap[dao.Article, domain.Article](arts, func(idx int, src dao.Article) domain.Article {
		return repo.toDomain(src)
	})
	return rs, err
}

func (repo *articleRepository) Sync(ctx context.Context, art domain.Article) (int, error) {
	return repo.dao.Sync(ctx, repo.toEntity(art))
}

func (repo *articleRepository) Update(ctx context.Context, art domain.Article) (int, error) {
	return repo.dao.UpdateById(ctx, repo.toEntity(art))
}

func (repo *articleRepository) Create(ctx context.Context, art domain.Article) (int, error) {
	return repo.dao.Create(ctx, repo.toEntity(art))
}

func (repo articleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		Status:   uint(art.Status),
		AuthorId: art.Author.Id,
		Ctime:    art.Ctime,
		Utime:    art.Utime,
	}
}

func (repo *articleRepository) toDomain(art dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Status:  domain.ArticleStatus(art.Status),
		Author: domain.Author{
			Id: art.AuthorId,
		},
		Ctime: art.Ctime,
		Utime: art.Utime,
	}
}
