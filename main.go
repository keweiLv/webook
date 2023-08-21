package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/keweiLv/webook/env"
	"github.com/keweiLv/webook/internal/repository"
	"github.com/keweiLv/webook/internal/repository/dao"
	"github.com/keweiLv/webook/internal/service"
	"github.com/keweiLv/webook/internal/web"
	"github.com/keweiLv/webook/internal/web/middleware"
	"github.com/keweiLv/webook/pkg/ginx/middlewares/ratelimit"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	_ "log"
	"net/http"
	_ "net/http"
	"time"
)

var v *viper.Viper

func main() {
	//v = initConfig()
	db := initDB()
	server := initWebServer()
	u := initUser(db)
	u.RegisterRoutes(server)
	//server := gin.Default()
	server.GET("hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "yes，it's ok")
	})
	server.Run(":8081")
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		println("这是第一个 middleware")
	})

	server.Use(func(ctx *gin.Context) {
		println("这是第二个 middleware")
	})

	// redis 限流
	redisClient := redis.NewClient(&redis.Options{
		Addr: env.Config.Redis.Addr,
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	server.Use(cors.New(cors.Config{
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
	}))
	//store := cookie.NewStore([]byte("secret"))

	store := memstore.NewStore([]byte("ZfLRVcUjaIjbJBiZtvMhgxNSqZT2nMFl"), []byte("0pCCmjIcuL4xE8WXY57fucEh6RpE88rt"))

	server.Use(sessions.Sessions("ssid", store))

	//server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePaths("/users/signup").IgnorePaths("/users/login").Build())

	// jwt
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().IgnorePaths("/users/signup").IgnorePaths("/users/login").Build())

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
	//username := v.GetString("mysql.username")
	//password := v.GetString("mysql.password")
	//host := v.GetString("mysql.host")
	//port := v.GetString("mysql.port")
	//dbname := v.GetString("mysql.dbname")
	//viperDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(env.Config.DB.DSN))
	//db, err := gorm.Open(mysql.Open(viperDsn))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initConfig() *viper.Viper {
	v := viper.New()                  // 添加配置文件搜索路径，点号为当前目录
	v.AddConfigPath("./remoteConf")   // 添加多个搜索目录
	v.SetConfigType("yml")            // 如果配置文件没有后缀，可以不用配置
	v.SetConfigName("dev_config.yml") // 文件名，没有后缀
	// 读取配置文件
	if err := v.ReadInConfig(); err == nil {
		log.Printf("use config file -> %s\n", v.ConfigFileUsed())
	} else {
		panic(err)
	}
	return v
}
