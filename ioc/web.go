package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/keweiLv/webook/internal/web"
	"github.com/keweiLv/webook/internal/web/middleware"
	"github.com/keweiLv/webook/pkg/ginx/middlewares/ratelimit"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitWebService(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		middleware.NewLoginJWTMiddlewareBuilder().
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").Build(),
		ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 是否允许你带 cookie 之类的东西
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowOrigins:     []string{"http://localhost:3000"},
		//AllowMethods: []string{"POST", "GET"},

		// 允许前端读取的返回 header
		ExposeHeaders: []string{"x-jwt-token"},

		// 开发环境
		//AllowOriginFunc: func(origin string) bool {
		//	if strings.HasPrefix(origin, "http://localhost") {
		//		// 你的开发环境
		//		return true
		//	}
		//	return strings.Contains(origin, "yourcompany.com")
		//},
		MaxAge: 12 * time.Hour,
	})
}
