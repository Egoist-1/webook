//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"start/webook/internal/repository"
	"start/webook/internal/repository/cache"
	"start/webook/internal/repository/dao"
	"start/webook/internal/repository/dao/article"
	"start/webook/internal/service"
	"start/webook/internal/service/sms/memory"
	"start/webook/internal/web"
)

var thirdProvider = wire.NewSet()

//go:generate wire ./
func InitUserTest(dao dao.UserDAO, redis redis.Cmdable) *web.UserHandle {
	wire.Build(
		//user
		repository.NewUserRepo,
		service.NewUserServiceImpl,
		web.NewUserHandle,
		//code
		service.NewCodeService,
		repository.NewCodeRepo,
		cache.NewCodeCacheRedis,
		//sms
		memory.NewMemory,
	)
	return new(web.UserHandle)
}

//go:generate wire ./
func InitArticleTest(db *gorm.DB) *web.ArticleHandle {
	wire.Build(
		web.NewArticleHandle,
		service.NewArticleService,
		repository.NewArticleRepository,
		article.NewArticleGormDao,
	)
	return new(web.ArticleHandle)
}
