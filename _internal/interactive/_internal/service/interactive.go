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
	Collected(ctx context.Context, biz string, aid int64, uid int64, collectionID int64, Collected bool) error
	Liked(ctx context.Context, biz string, aid int64, uid int64, liked bool) error
	ReadHistory(ctx context.Context, biz string, bizId int64, uid int64) error
}

func NewInterService(repo repository.InteractiveRepository) InteractiveService {
	return &IntrService{repo: repo}
}

type IntrService struct {
	repo repository.InteractiveRepository
}

func (i *IntrService) ReadHistory(ctx context.Context, biz string, bizId int64, uid int64) error {
	return i.repo.ReadHistory(ctx, biz, bizId, uid)
}

func (i *IntrService) Collected(ctx context.Context, biz string, aid int64, uid int64, collectionID int64, liked bool) error {
	return i.repo.Collected(ctx, biz, aid, uid, collectionID, liked)
}

func (i *IntrService) Liked(ctx context.Context, biz string, aid int64, uid int64, liked bool) error {
	return i.repo.Liked(ctx, biz, aid, uid, liked)
}

func (i *IntrService) TopN(ctx context.Context, now time.Time, batch int) ([]domain.Interactive, error) {
	//TODO implement me
	panic("implement me")
}

func (i *IntrService) GetIntr(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interactive, error) {
	return i.repo.GetIntr(ctx, biz, bizId, uid)
}
