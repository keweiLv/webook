package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/keweiLv/webook/internal/repository"
	"github.com/keweiLv/webook/internal/repository/dao"
	"github.com/keweiLv/webook/internal/service"
	"github.com/keweiLv/webook/internal/web"
	"github.com/keweiLv/webook/internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	_ "log"
	_ "net/http"
	"time"
)

func main() {
	db := initDB()
	server := initWebServer()
	u := initUser(db)
	u.RegisterRoutes(server)
	server.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		println("这是第一个 middleware")
	})

	server.Use(func(ctx *gin.Context) {
		println("这是第二个 middleware")
	})
	server.Use(cors.New(cors.Config{
		// 是否允许你带 cookie 之类的东西
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowOrigins:     []string{"http://localhost:3000"},
		//AllowMethods: []string{"POST", "GET"},

		//ExposeHeaders: []string{"x-jwt-token"},

		// 开发环境
		//AllowOriginFunc: func(origin string) bool {
		//	if strings.HasPrefix(origin, "http://localhost") {
		//		// 你的开发环境
		//		return true
		//	}
		//	return strings.Contains(origin, "yourcompany.com")
		//},
		MaxAge: 12 * time.Hour,
	}))
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("ssid", store))

	server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePaths("/users/signup").IgnorePaths("/users/login").Build())

	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/webook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
