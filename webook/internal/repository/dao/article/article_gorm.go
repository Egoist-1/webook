package article

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"start/webook/internal/e"
	"time"
)

func NewArticleGormDao(db *gorm.DB) ArticleDao {
	return &articleGormDao{
		db: db,
	}
}

type articleGormDao struct {
	db *gorm.DB
}

func (dao *articleGormDao) PubFinById(ctx context.Context, uid int, aid int, status uint) (Article, error) {
	var art = Article{}
	err := dao.db.WithContext(ctx).Model(&Article{}).
		Where("id = ? AND author_id = ?, status = ?", uid, aid, status).
		First(&art).Error
	return art, err
}

func (dao *articleGormDao) FindById(ctx context.Context, uid int, aid int) (Article, error) {
	var art Article
	err := dao.db.WithContext(ctx).Model(&Article{}).
		Where("id = ? AND author_id = ?", aid, uid).
		First(&art).
		Error
	return art, err
}

func (dao *articleGormDao) GetUnpublishList(ctx context.Context, uid int, limit int, offset int) (arts []Article, err error) {
	err = dao.db.WithContext(ctx).
		Model(&Article{}).
		Select("id,time,status,author_id,ctime,utime").
		Where("Author_id = ?", uid).
		Limit(limit).
		Offset(offset).
		Order("Ctime Desc").
		Find(&arts).Error
	return arts, err
}

func (dao *articleGormDao) Sync(ctx context.Context, art Article) (int, error) {
	var (
		//制作库的Id
		id int
	)
	err := dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		txdao := NewArticleGormDao(tx)
		if art.Id > 0 {
			id, err = txdao.UpdateById(ctx, art)
		} else {
			id, err = txdao.Create(ctx, art)
		}
		if err != nil {
			return err
		}
		publish := ArticlePublish(art)
		now := time.Now().UnixMilli()
		err = tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":   publish.Title,
				"content": publish.Content,
				"status":  publish.Status,
				"utime":   now,
			}),
		}).Create(&publish).Error
		return err
	})
	return id, err
}

func (dao *articleGormDao) UpdateById(ctx context.Context, art Article) (int, error) {
	art.Utime = time.Now().UnixMilli()
	result := dao.db.WithContext(ctx).
		Where("id=? AND author_id = ? ", art.Id, art.AuthorId).
		Updates(&art)
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		zap.L().Warn("articleGormDao 请求参数错误", zap.Any("请求", art))
		return 0, e.NewServerErr("articleGormDao 请求参数错误", "")
	}
	return art.Id, result.Error
}

func (dao *articleGormDao) Create(ctx context.Context, art Article) (int, error) {
	now := time.Now().UnixMilli()
	art.Utime = now
	art.Ctime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}
