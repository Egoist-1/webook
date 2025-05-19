package service

import (
	"context"
	"time"
	"webook/_internal/interactive/_internal/domain"
	"webook/_internal/interactive/_internal/repository"
)

//go:generate mockgen -source=interactive.go -package=svcmocks -destination=mocks/interactive.mock.go InteractiveService
type InteractiveService interface {
	GetIntr(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interactive, error)
	TopN(ctx context.Context, now time.Time, batch int) ([]domain.Interactive, error)
	Collected(ctx context.Context, biz string, aid int64, uid int64, collectionID int64) error
	Liked(ctx context.Context, biz string, aid int64, uid int64) error
	ReadHistory(ctx context.Context, biz string, bizId int64, uid int64) error
	GetReadHistory(ctx context.Context, biz string, uid int64, offset, limit int) ([]domain.ReadHistory, error)
	CreateCollection(ctx context.Context, biz string, uid int64, cname string) error
	CancelCollection(ctx context.Context, biz string, uid int64, aid, cid int64) error
	GetCollectionList(ctx context.Context, uid int64, biz string) ([]domain.Collection, error)
	CollectionDetail(ctx context.Context, biz string, cid int64) (bizIds []int64, err error)
	CancelLike(ctx context.Context, biz string, uid int64, aid int64) error
}

func NewInterService(repo repository.InteractiveRepository) InteractiveService {
	return &IntrService{repo: repo}
}

type IntrService struct {
	repo repository.InteractiveRepository
}

func (i *IntrService) CancelLike(ctx context.Context, biz string, uid int64, aid int64) error {
	return i.repo.CancelLike(ctx, biz, uid, aid)
}

func (i *IntrService) CollectionDetail(ctx context.Context, biz string, cid int64) (bizIds []int64, err error) {
	return i.repo.CollectionDetail(ctx, biz, cid)
}

func (i *IntrService) GetCollectionList(ctx context.Context, uid int64, biz string) ([]domain.Collection, error) {
	return i.repo.GetCollectionList(ctx, uid, biz)
}

func (i *IntrService) CancelCollection(ctx context.Context, biz string, uid int64, aid, cid int64) error {
	return i.repo.CancelCollection(ctx, biz, uid, aid, cid)
}

func (i *IntrService) CreateCollection(ctx context.Context, biz string, uid int64, cname string) error {
	return i.repo.CreateCollection(ctx, biz, uid, cname)
}

func (i *IntrService) GetReadHistory(ctx context.Context, biz string, uid int64, offset, limit int) ([]domain.ReadHistory, error) {
	return i.repo.GetReadHistory(ctx, biz, uid, offset, limit)
}

func (i *IntrService) ReadHistory(ctx context.Context, biz string, bizId int64, uid int64) error {
	return i.repo.ReadHistory(ctx, biz, bizId, uid)
}

func (i *IntrService) Collected(ctx context.Context, biz string, aid int64, uid int64, collectionID int64) error {
	return i.repo.Collected(ctx, biz, aid, uid, collectionID)
}

func (i *IntrService) Liked(ctx context.Context, biz string, aid int64, uid int64) error {
	return i.repo.Liked(ctx, biz, aid, uid)
}

func (i *IntrService) TopN(ctx context.Context, now time.Time, batch int) ([]domain.Interactive, error) {
	//TODO implement me
	panic("implement me")
}

func (i *IntrService) GetIntr(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interactive, error) {
	return i.repo.GetIntr(ctx, biz, bizId, uid)
}
