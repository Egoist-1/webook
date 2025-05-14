package service

import (
	"context"
	"start/webook/interactive/_internal/repository"
	"start/webook/interactive/internal/domain"
	"time"
)

//go:generate mockgen -source=interactive.go -package=svcmocks -destination=mocks/interactive.mock.go InteractiveService
type InteractiveService interface {
	GetIntr(ctx context.Context, biz string, bizId int64, uid int) (domain.Interactive, error)
	TopN(ctx context.Context, now time.Time, batch int) ([]domain.Interactive, error)
}

func NewIntrService(repo repository.InteractiveRepository) InteractiveService {
	return &IntrService{repo: repo}
}

type IntrService struct {
	repo repository.InteractiveRepository
}

func (i *IntrService) TopN(ctx context.Context, now time.Time, batch int) ([]domain.Interactive, error) {
	//TODO implement me
	panic("implement me")
}

func (i *IntrService) GetIntr(ctx context.Context, biz string, bizId int64, uid int) (domain.Interactive, error) {
	return i.repo.GetIntr(ctx, biz, bizId, uid)
}
