package startup

import (
	"context"
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"start/webook/internal/repository/dao"
	"time"
)

var db *gorm.DB

func InitTestDB() *gorm.DB {
	if db == nil {
		dsn := "root:root@tcp(localhost:3306)/webook"
		sqlDB, err := sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}
		for {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err = sqlDB.PingContext(ctx)
			cancel()
			if err == nil {
				break
			}
			log.Println("等待连接 MySQL", err)
		}
		db, err = gorm.Open(mysql.Open(dsn))
		if err != nil {
			panic(err)
		}
		err = dao.InitTables(db)
		if err != nil {
			panic(err)
		}
		//db = db.Debug()
	}
	return db
}
