// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/keweiLv/webook/internal/repository"
	"github.com/keweiLv/webook/internal/repository/cache"
	"github.com/keweiLv/webook/internal/repository/dao"
	"github.com/keweiLv/webook/internal/service"
	"github.com/keweiLv/webook/internal/web"
	"github.com/keweiLv/webook/ioc"
)

import (
	_ "log"
	_ "net/http"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := ioc.InitRedis()
	v2 := ioc.InitMiddlewares(cmdable)
	db := ioc.InitDB()
	userDAO := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	codeRedisCache := cache.NewCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeRedisCache)
	smsService := ioc.InitSMSService()
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService)
	engine := ioc.InitWebService(v2, userHandler)
	return engine
}
