package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"token_surveillance_system/extend/code"
	"token_surveillance_system/extend/jwt"
	"token_surveillance_system/extend/redis"
	"token_surveillance_system/extend/utils"
)

//jwt认证中间件

func JWTAuth() gin.HandlerFunc  {
	return func(context *gin.Context) {
		//获取token
		token := context.Request.Header.Get("Authorization")
		if token == ""{
			utils.ResponseFormat(context,code.TokenNotFound,nil)
			context.Abort()
			return
		}
		//解析token
		jwtInstance := jwt.NewJWT()
		claims,err := jwtInstance.ParseToken(token)
		if err!=nil{
			utils.ResponseFormat(context,code.TokenInvalid,nil)
			context.Abort()
			return
		}
		//获取缓存中的token信息
		tokenCache,err := redis.Get("TOKEN:"+claims.Email)
		if err != nil {
			log.Error().Msgf("jwt auth redis get :%v",err.Error())
			utils.ResponseFormat(context,code.ServiceInsideError,nil)
			context.Abort()
			return
		}
		//用户注销或者token失效
		if tokenCache != token {
			log.Error().Msg("user singout or token invaild")
			utils.ResponseFormat(context,code.TokenInvalid,nil)
			context.Abort()
			return
		}
		context.Set("claims",claims)
		context.Next()
	}
}
