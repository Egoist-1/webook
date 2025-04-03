package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type InteractiveDao interface {
	FindByBizId(ctx context.Context, biz string, aid int64, uid int) (Interactive, error)
	LikedInfo(ctx context.Context, biz string, id int64, uid int) (UserIntrInfo, error)
	IncrReadCnt(ctx context.Context, aid int) error
}

func NewIntrDao(db *gorm.DB) InteractiveDao {
	return &intrDao{
		db: db,
	}
}

type intrDao struct {
	db *gorm.DB
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

func (dao *intrDao) FindByBizId(ctx context.Context, biz string, aid int64, uid int) (Interactive, error) {
	var interactive Interactive
	err := dao.db.WithContext(ctx).
		Model(&Interactive{}).
		Where("biz = ? AND biz_id = ?", biz, aid).
		First(&interactive).
		Error
	return interactive, err
}

func (dao *intrDao) LikedInfo(ctx context.Context, biz string, id int64, uid int) (UserIntrInfo, error) {
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
	Id         int
	Biz        string
	BizId      int64
	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Ctime      int64
	Utime      int64
}

type UserIntrInfo struct {
	Id        int
	Uid       int
	Biz       string
	BizId     int64
	Liked     bool
	Collected bool
	//软删除
	Status int
	Ctime  int64
	Utime  int64
}
