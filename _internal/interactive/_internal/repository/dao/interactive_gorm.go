package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type InteractiveDao interface {
	FindByBizId(ctx context.Context, biz string, aid int64, uid int64) (Interactive, error)
	LikedInfo(ctx context.Context, biz string, id int64, uid int64) (UserIntrInfo, error)
	IncrReadCnt(ctx context.Context, aid int) error
	Collected(ctx context.Context, biz string, aid int64, uid int64, collectionID int64, liked bool) error
	Liked(ctx context.Context, biz string, aid int64, uid int64, Collected bool) error
	ReadHistory(ctx context.Context, biz string, bizId int64, uid int64) error
}

func NewIntrDao(db *gorm.DB) InteractiveDao {
	return &intrDao{
		db: db,
	}
}

type intrDao struct {
	db *gorm.DB
}

func (dao *intrDao) ReadHistory(ctx context.Context, biz string, bizId int64, uid int64) error {
	now := time.Now().UnixMilli()
	err := dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"utime": time.Now(),
		}),
	}).Create(&ReadHistory{
		Uid:   uid,
		Biz:   biz,
		BizId: bizId,
		Ctime: now,
		Utime: now,
	}).Error
	return err
}

func (dao *intrDao) Collected(ctx context.Context, biz string, aid int64, uid int64, cid int64, collected bool) error {
	var expr string
	if collected {
		expr = "collected_cnt + 1"
	} else {
		expr = "collected_cnt - 1"
	}
	now := time.Now().UnixMilli()
	err := dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Interactive{}).Where("biz_id = ?", aid).
			Update("collected_cnt", gorm.Expr(expr)).Error
		if err != nil {
			return err
		}
		err = tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"collected": true,
			}),
		}).Create(&Collection{
			Uid:            uid,
			Biz:            biz,
			BizId:          aid,
			CollectionId:   uid,
			CollectionName: "默认收藏夹",
			Collected:      collected,
			Ctime:          now,
			Utime:          now,
		}).Error
		return err
	})
	return err
}

func (dao *intrDao) Liked(ctx context.Context, biz string, aid int64, uid int64, liked bool) error {
	var expr string
	if liked {
		expr = "liked + 1"
	} else {
		expr = "liked - 1"
	}
	//在这里同时增加like的计数和创建用户关联
	now := time.Now().UnixMilli()
	err := dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//增加点赞记录
		err := tx.Model(&Interactive{}).Where("biz_id = ?", aid).
			Update("like_cnt", gorm.Expr(expr)).Error
		//增加用户点赞记录
		err = tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"liked": true,
			}),
		}).Create(&UserIntrInfo{
			Uid:   uid,
			Biz:   biz,
			BizId: aid,
			Liked: liked,
			Ctime: now,
			Utime: now,
		}).Error
		return err
	})
	return err
}

func (dao *intrDao) IncrReadCnt(ctx context.Context, aid int) error {
	now := time.Now().UnixMilli()
	err := dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"read_cnt": gorm.Expr("read_cnt + 1"),
		}),
	}).Create(&Interactive{
		Biz:     "article",
		BizId:   int64(aid),
		ReadCnt: 1,
		Ctime:   now,
		Utime:   now,
	}).Error
	return err
}

func (dao *intrDao) FindByBizId(ctx context.Context, biz string, aid int64, uid int64) (Interactive, error) {
	var interactive Interactive
	err := dao.db.WithContext(ctx).
		Model(&Interactive{}).
		Where("biz = ? AND biz_id = ?", biz, aid).
		First(&interactive).
		Error
	return interactive, err
}

func (dao *intrDao) LikedInfo(ctx context.Context, biz string, id int64, uid int64) (UserIntrInfo, error) {
	var ulInfo UserIntrInfo
	err := dao.db.WithContext(ctx).Model(&UserIntrInfo{}).
		Where("biz = ? and biz_id = ? AND uid = ?", biz, id, uid).
		First(&ulInfo).Error
	if err == gorm.ErrRecordNotFound {
		return ulInfo, nil
	}
	return ulInfo, err
}

type Interactive struct {
	Id  int64
	Biz string
	//也就是 文章的id
	BizId      int64
	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Ctime      int64
	Utime      int64
}

// TODO biz bizID uid 添加唯一索引
type UserIntrInfo struct {
	Id        int64
	Uid       int64
	Biz       string
	BizId     int64
	Liked     bool
	Collected bool
	//软删除
	Status int
	Ctime  int64
	Utime  int64
}

// TODO biz bizID uid,collectionId 添加唯一索引
type Collection struct {
	Id    int64
	Uid   int64
	Biz   string
	BizId int64
	//文件夹id
	CollectionId int64
	//文件夹名称
	CollectionName string
	Collected      bool
	Ctime          int64
	Utime          int64
}
type ReadHistory struct {
	Id    int64
	Uid   int64
	Biz   string
	BizId int64
	Ctime int64
	Utime int64
}
