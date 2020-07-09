package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	v1 "token_surveillance_system/controller/v1"
	"token_surveillance_system/extend/conf"
	"token_surveillance_system/middleware"
	"token_surveillance_system/service"
)

func InitRouter()*gin.Engine{
	r:=gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(conf.ServerConf.RunMode)
	//设置跨域
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: conf.CORSConf.AllowAllOrigins,
		AllowMethods: conf.CORSConf.AllowMethods,
		AllowHeaders: conf.CORSConf.AllowHeaders,
		ExposeHeaders: conf.CORSConf.ExposeHeaders,
		AllowCredentials: conf.CORSConf.AllowCredentials,
		MaxAge: conf.CORSConf.MaxAge * time.Hour,
	}))
	uploadService := service.UploadService{}
	r.StaticFS("/upload/img",http.Dir(uploadService.GetImgFullPath()))
	apiV1 := r.Group("/api/v1")
	authController := v1.AuthController{}
	{
		apiV1.POST("/auth/singn",authController.Signup)//注册
		apiV1.POST("/auth/signin",authController.Signin)//登录
		userController:=new(v1.UserController)
		apiV1.Use(middleware.JWTAuth())
		{
			apiV1.POST("/auth/signout",authController.Signout)//注销
			apiV1.GET("/user",userController.FindUser)//查看用户信息
			apiV1.PATCH("/user/pass",userController.Alterpass)
		}
	}
}
