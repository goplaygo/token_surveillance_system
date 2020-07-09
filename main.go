package main

import (
	"fmt"
	"token_surveillance_system/extend/conf"
	"token_surveillance_system/extend/logger"
	"token_surveillance_system/extend/redis"
	"token_surveillance_system/extend/validator"
	"token_surveillance_system/models"
	"token_surveillance_system/router"
)

func main(){
	conf.Setup()
	logger.Setup()
	redis.Setup()
	models.Setup()
	validator.Setup()
	router := router.InitRouter()
	router.Run(fmt.Sprintf(":%d",conf.ServerConf.Port))
}
