package integration

import (
	"context"
	"fmt"
	"start/webook/internal/integration/startup"
	"start/webook/internal/repository/dao"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	db := startup.InitTestDB()
	go func() {
		var u dao.User
		tx := db.Begin()
		_ = tx.Model(&dao.User{}).Where("id = ? ", 1).Update("name", "13").Error
		time.Sleep(time.Second * 10)
		tx.Commit()
	}()
	time.Sleep(time.Second * 2)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	fmt.Println("开始执行")
	err := db.WithContext(ctx).Model(&dao.User{}).Where("id = ? ", 1).Update("name", "23").Error
	fmt.Println("结束")
	fmt.Println(err)
	time.Sleep(time.Second * 10)
}
