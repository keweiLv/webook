package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	_ "log"
	"net/http"
	_ "net/http"
)

var v *viper.Viper

func main() {
	//v = initConfig()
	//db := initDB()
	//server := initWebServer()
	//rdb := initRedis()
	//u := initUser(db, rdb)
	//u.RegisterRoutes(server)

	//wire 实现
	server := InitWebServer()

	//server := gin.Default()
	server.GET("hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "yes，it's ok")
	})
	server.Run(":8081")
}
