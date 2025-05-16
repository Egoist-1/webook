package service

import (
	"context"
	"start/webook/interactive/_internal/domain"
	"start/webook/interactive/_internal/repository"
	"time"
)

//go:generate mockgen -source=interactive.go -package=svcmocks -destination=mocks/interactive.mock.go InteractiveService
type InteractiveService interface {
	GetIntr(ctx context.Context, biz string, bizId int64, uid int) (domain.Interactive, error)
	TopN(ctx context.Context, now time.Time, batch int) ([]domain.Interactive, error)
	Collected(ctx context.Context, biz string, aid int64, uid int64, collectionID int64) error
	Liked(ctx context.Context, biz string, aid int64, uid int64) error
	CancelTheLike(ctx context.Context, biz string, aid int64, uid int64) error
}

func NewInterService(repo repository.InteractiveRepository) InteractiveService {
	return &IntrService{repo: repo}
}

type IntrService struct {
	repo repository.InteractiveRepository
}

func (i *IntrService) CancelTheLike(ctx context.Context, biz string, aid int64, uid int64) error {
	//TODO implement me
	panic("implement me")
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

func (i *IntrService) GetIntr(ctx context.Context, biz string, bizId int64, uid int) (domain.Interactive, error) {
	return i.repo.GetIntr(ctx, biz, bizId, uid)
}
