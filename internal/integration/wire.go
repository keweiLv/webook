//go:build wireinject

package integration

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/keweiLv/webook/internal/repository"
	"github.com/keweiLv/webook/internal/repository/cache"
	"github.com/keweiLv/webook/internal/repository/dao"
	"github.com/keweiLv/webook/internal/service"
	"github.com/keweiLv/webook/internal/web"
	"github.com/keweiLv/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,

		dao.NewUserDAO,

		cache.NewUserCache,
		cache.NewCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,

		service.NewUserService,
		service.NewCodeService,

		ioc.InitSMSService,
		web.NewUserHandler,

		ioc.InitWebService,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
