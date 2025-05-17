package ioc

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
	"webook/_internal/article/_internal/repository/dao"
)

func InitGorm() *gorm.DB {
	db, err := gorm.Open(mysql.Open(viper.GetString("db.dsn")), &gorm.Config{
		Logger: logger.New(gormLogger(zap.L().Debug), logger.Config{
			SlowThreshold:             time.Millisecond * 100,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
			LogLevel:                  logger.Info,
		}),
	})
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

type gormLogger func(string, ...zap.Field)

func (g gormLogger) Printf(s string, i ...interface{}) {
	g(s, zap.Any(s, i))
}
