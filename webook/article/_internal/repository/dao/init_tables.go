package dao

import (
	"gorm.io/gorm"
	dao2 "start/webook/user/_internal/repository/dao"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&dao2.User{},
		Article{},
		ArticlePublish{},
	)
}
