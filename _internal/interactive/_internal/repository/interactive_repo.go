package repository

import (
	"context"
	"golang.org/x/sync/errgroup"
	"webook/_internal/interactive/_internal/domain"
	"webook/_internal/interactive/_internal/repository/dao"
)

type InteractiveRepository interface {
	GetIntr(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error)
	IncrReadCnt(ctx context.Context, aid int) error

	Collected(ctx context.Context, biz string, aid int64, uid int64, collectionID int64, liked bool) error
	Liked(ctx context.Context, biz string, aid int64, uid int64, Collected bool) error
	ReadHistory(ctx context.Context, biz string, bizId int64, uid int64) error
}

func NewIntrRepo(dao dao.InteractiveDao) InteractiveRepository {
	return &IntrRepo{dao: dao}
}

type IntrRepo struct {
	dao dao.InteractiveDao
}

func (repo *IntrRepo) ReadHistory(ctx context.Context, biz string, bizId, uid int64) error {
	return repo.dao.ReadHistory(ctx, biz, bizId, uid)
}

func (repo *IntrRepo) Collected(ctx context.Context, biz string, aid int64, uid int64, cid int64, Collected bool) error {
	return repo.dao.Collected(ctx, biz, aid, uid, cid, Collected)
}

func (repo *IntrRepo) Liked(ctx context.Context, biz string, aid int64, uid int64, liked bool) error {
	return repo.dao.Liked(ctx, biz, aid, uid, liked)
}
func (repo *IntrRepo) IncrReadCnt(ctx context.Context, aid int) error {
	return repo.dao.IncrReadCnt(ctx, aid)
}

func (repo *IntrRepo) GetIntr(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error) {
	var (
		eg      errgroup.Group
		intrDao dao.Interactive
		ulInfo  dao.UserIntrInfo
		err     error
	)
	eg.Go(func() error {
		intrDao, err = repo.dao.FindByBizId(ctx, biz, id, uid)
		return err
	})
	eg.Go(func() error {
		ulInfo, err = repo.dao.LikedInfo(ctx, biz, id, uid)
		return err
	})
	err = eg.Wait()
	result := repo.toDomain(intrDao, ulInfo.Liked, ulInfo.Collected)
	return result, err
}
func (repo *IntrRepo) toDomain(intr dao.Interactive, liked, Collected bool) domain.Interactive {
	return domain.Interactive{
		Biz:        intr.Biz,
		BizId:      intr.BizId,
		ReadCnt:    intr.ReadCnt,
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		Liked:      liked,
		Collected:  Collected,
		Ctime:      intr.Ctime,
		Utime:      intr.Utime,
	}
}

func (repo *IntrRepo) toEntity(intr domain.Interactive) dao.Interactive {
	return dao.Interactive{
		Biz:        intr.Biz,
		BizId:      intr.BizId,
		ReadCnt:    intr.ReadCnt,
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		Ctime:      intr.Ctime,
		Utime:      intr.Utime,
	}
}
