package repository

import (
	"context"
	"golang.org/x/sync/errgroup"
	"webook/_internal/interactive/_internal/domain"
	"webook/_internal/interactive/_internal/repository/dao"
	"webook/pkg/tools/slicex"
)

type InteractiveRepository interface {
	GetIntr(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error)
	IncrReadCnt(ctx context.Context, aid int) error

	Collected(ctx context.Context, biz string, aid int64, uid int64, collectionID int64) error
	Liked(ctx context.Context, biz string, aid int64, uid int64) error
	ReadHistory(ctx context.Context, biz string, bizId int64, uid int64) error
	GetReadHistory(ctx context.Context, biz string, uid int64, offset int, limit int) ([]domain.ReadHistory, error)
	CreateCollection(ctx context.Context, biz string, uid int64, cname string) error
	CancelCollection(ctx context.Context, biz string, uid int64, aid int64, cid int64) error
	GetCollectionList(ctx context.Context, uid int64, biz string) ([]domain.Collection, error)
	CollectionDetail(ctx context.Context, biz string, cid int64) ([]int64, error)
	CancelLike(ctx context.Context, biz string, uid int64, aid int64) error
}

func NewIntrRepo(dao dao.InteractiveDao) InteractiveRepository {
	return &IntrRepo{dao: dao}
}

type IntrRepo struct {
	dao dao.InteractiveDao
}

func (repo *IntrRepo) CancelLike(ctx context.Context, biz string, uid int64, aid int64) error {
	return repo.dao.CancelLike(ctx, biz, uid, aid)
}

func (repo *IntrRepo) CollectionDetail(ctx context.Context, biz string, cid int64) ([]int64, error) {
	return repo.dao.CollectionDetail(ctx, biz, cid)
}

func (repo *IntrRepo) GetCollectionList(ctx context.Context, uid int64, biz string) ([]domain.Collection, error) {
	list, err := repo.dao.GetCollectionList(ctx, uid, biz)
	if err != nil {
		return nil, err
	}
	res := slicex.SliceMap[dao.Collection, domain.Collection](list, func(idx int, src dao.Collection) domain.Collection {
		return repo.CollectionToDomain(src)
	})
	return res, nil
}

func (repo *IntrRepo) CancelCollection(ctx context.Context, biz string, uid int64, aid int64, cid int64) error {
	return repo.dao.CancelCollection(ctx, biz, uid, aid, cid)
}

func (repo *IntrRepo) CreateCollection(ctx context.Context, biz string, uid int64, cname string) error {
	return repo.dao.CreateCollection(ctx, biz, uid, cname)
}

func (repo *IntrRepo) GetReadHistory(ctx context.Context, biz string, uid int64, offset int, limit int) ([]domain.ReadHistory, error) {
	history, err := repo.dao.GetReadHistory(ctx, biz, uid, offset, limit)
	if err != nil {
		return nil, err
	}
	domains := slicex.SliceMap[dao.ReadHistory, domain.ReadHistory](history, func(idx int, src dao.ReadHistory) domain.ReadHistory {
		return repo.readHistoryToDomain(src)
	})
	return domains, err
}

func (repo *IntrRepo) ReadHistory(ctx context.Context, biz string, bizId, uid int64) error {
	return repo.dao.ReadHistory(ctx, biz, bizId, uid)
}

func (repo *IntrRepo) Collected(ctx context.Context, biz string, aid int64, uid int64, cid int64) error {
	return repo.dao.Collected(ctx, biz, aid, uid, cid)
}

func (repo *IntrRepo) Liked(ctx context.Context, biz string, aid int64, uid int64) error {
	return repo.dao.Liked(ctx, biz, aid, uid)
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
func (repo *IntrRepo) readHistoryToDomain(r dao.ReadHistory) domain.ReadHistory {
	return domain.ReadHistory{
		Id:    r.Id,
		Uid:   r.Uid,
		Biz:   r.Biz,
		BizId: r.BizId,
		Ctime: r.Ctime,
		Utime: r.Utime,
	}
}
func (repo *IntrRepo) readHistoryToEntity(r domain.ReadHistory) dao.ReadHistory {
	return dao.ReadHistory{
		Id:    r.Id,
		Uid:   r.Uid,
		Biz:   r.Biz,
		BizId: r.BizId,
		Ctime: r.Ctime,
		Utime: r.Utime,
	}
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
func (repo *IntrRepo) CollectionToDomain(intr dao.Collection) domain.Collection {
	return domain.Collection{
		Id:             intr.Id,
		Biz:            intr.Biz,
		Uid:            intr.Uid,
		CollectionName: intr.CollectionName,
		Ctime:          intr.Ctime,
		Utime:          intr.Utime,
	}
}

func (repo *IntrRepo) CollectionToEntity(intr domain.Collection) dao.Collection {
	return dao.Collection{
		Id:             intr.Id,
		Biz:            intr.Biz,
		Uid:            intr.Uid,
		CollectionName: intr.CollectionName,
		Ctime:          intr.Ctime,
		Utime:          intr.Utime,
	}
}
