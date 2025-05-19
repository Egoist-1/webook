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
	Collected(ctx context.Context, biz string, aid int64, uid int64, collectionID int64) error
	Liked(ctx context.Context, biz string, aid int64, uid int64) error
	ReadHistory(ctx context.Context, biz string, bizId int64, uid int64) error
	GetReadHistory(ctx context.Context, biz string, uid int64, offset int, limit int) ([]ReadHistory, error)
	CreateCollection(ctx context.Context, biz string, uid int64, cname string) error
	CancelCollection(ctx context.Context, biz string, uid int64, aid int64, cid int64) error
	GetCollectionList(ctx context.Context, uid int64, biz string) ([]Collection, error)
	CollectionDetail(ctx context.Context, biz string, cid int64) ([]int64, error)
	CancelLike(ctx context.Context, biz string, uid int64, aid int64) error
}

func NewIntrDao(db *gorm.DB) InteractiveDao {
	return &intrDao{
		db: db,
	}
}

type intrDao struct {
	db *gorm.DB
}

func (dao *intrDao) CancelLike(ctx context.Context, biz string, uid int64, aid int64) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		return nil
	})
}

func (dao *intrDao) CollectionDetail(ctx context.Context, biz string, cid int64) ([]int64, error) {
	var ids []int64
	err := dao.db.WithContext(ctx).Model(&UserCollectionBiz{}).
		Select("biz_id").
		Where("biz = ? and cid = ?", biz).
		Find(&ids).Error
	return ids, err
}

func (dao *intrDao) GetCollectionList(ctx context.Context, uid int64, biz string) ([]Collection, error) {
	var list []Collection
	err := dao.db.WithContext(ctx).Model(&Collection{}).
		Where("biz = ? and uid = ?", biz, uid).
		Find(&[]Collection{}).Error
	return list, err
}

func (dao *intrDao) CancelCollection(ctx context.Context, biz string, uid int64, aid int64, cid int64) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Model(&Interactive{}).
			Where("biz = ? and uid = ? and cid = ?", biz, uid, cid).
			Update()
		err := tx.Model(&UserIntrInfo{}).
			Where("biz = ? and uid = ? aid = ?", biz, uid, aid).
			Update("collected = ?", false).
			Error
		if err != nil {
			return err
		}
		//删除收藏夹内的
		err = tx.Model(&UserCollectionBiz{}).Delete(&UserCollectionBiz{
			Id:    0,
			Cid:   cid,
			BizId: aid,
			Biz:   biz,
		}).Error
		return err
	})
}

func (dao *intrDao) CreateCollection(ctx context.Context, biz string, uid int64, cname string) error {
	now := time.Now().UnixMilli()
	err := dao.db.WithContext(ctx).Create(&Collection{
		Biz:            biz,
		Uid:            uid,
		CollectionName: cname,
		Ctime:          now,
		Utime:          now,
	}).Error
	return err
}

func (dao *intrDao) GetReadHistory(ctx context.Context, biz string, uid int64, offset int, limit int) ([]ReadHistory, error) {
	var history []ReadHistory
	err := dao.db.WithContext(ctx).Model(&ReadHistory{}).
		Where("biz = ? AND uid = ?", biz, uid).
		Order("utime desc").
		Offset(offset).
		Limit(limit).
		Find(&history).Error
	return history, err
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

func (dao *intrDao) Collected(ctx context.Context, biz string, aid int64, uid int64, cid int64) error {

	now := time.Now().UnixMilli()
	err := dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"collect_cnt": gorm.Expr("`collect_cnt`+1"),
				"utime":       now,
			}),
		}).Create(&Interactive{
			CollectCnt: 1,
			Ctime:      now,
			Utime:      now,
			Biz:        biz,
			BizId:      aid,
		}).Error
		err = tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"collected": true,
			}),
		}).Create(&UserIntrInfo{
			Uid:       uid,
			Biz:       biz,
			BizId:     aid,
			Collected: true,
			Ctime:     now,
			Utime:     now,
		}).Error
		if err != nil {
			return err
		}
		//创建
		err = tx.Model(&UserCollectionBiz{}).
			Create(&UserCollectionBiz{
				Cid:   cid,
				BizId: aid,
				Biz:   biz,
				Ctime: now,
				Utime: now,
			}).Error
		return err
	})
	return err
}

func (dao *intrDao) Liked(ctx context.Context, biz string, aid int64, uid int64) error {

	//在这里同时增加like的计数和创建用户关联
	now := time.Now().UnixMilli()
	err := dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//增加点赞记录
		err := tx.Model(&Interactive{}).Where("biz_id = ?", aid).
			Update("like_cnt", gorm.Expr("liked + 1")).Error
		//增加用户点赞记录
		err = tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"liked": true,
			}),
		}).Create(&UserIntrInfo{
			Uid:   uid,
			Biz:   biz,
			BizId: aid,
			Liked: true,
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

// 收藏夹
// TODO uid biz name 添加唯一索引
type Collection struct {
	Id  int64
	Biz string
	Uid int64
	//文件夹名称
	CollectionName string
	Ctime          int64
	Utime          int64
}

// 用户收藏作品
// TODO biz bizid cid uid 添加唯一索引
type UserCollectionBiz struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 收藏夹 ID
	// 作为关联关系中的外键，我们这里需要索引
	Cid   int64  `gorm:"index"`
	BizId int64  `gorm:"uniqueIndex:biz_type_id_uid"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:biz_type_id_uid"`
	Ctime int64
	Utime int64
}

// 阅读记录
type ReadHistory struct {
	Id    int64
	Uid   int64
	Biz   string
	BizId int64
	Ctime int64
	Utime int64
}
