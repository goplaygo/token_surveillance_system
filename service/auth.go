package service

import (
	"time"
	"token_surveillance_system/extend/conf"
	"token_surveillance_system/extend/jwt"
	"token_surveillance_system/extend/redis"
	"token_surveillance_system/models"
	goJWT"github.com/dgrijalva/jwt-go"
)

type AuthService struct {
	User *models.User
}

//根据 GenerateToken 生成token
func (au *AuthService)GenerateToken(user models.User) (string , error ){
	jwtInstance := jwt.NewJWT()
	nowTime := time.Now()
	expireTime := time.Duration(conf.ServerConf.JWTExpire)
	claims := jwt.CustomClaims{
		ID: user.ID,
		UserName:user.UserName,
		Email:user.Email,
		StandardClaims:goJWT.StandardClaims{
			ExpiresAt: nowTime.Add(expireTime * time.Hour).Unix(),
			Issuer: "monitor",
		},
	}
	token , err:= jwtInstance.CreatToken(claims)
	if err != nil {
		return "", err
	}
	//设置redis缓存
	const hourSecs int = 60*60
	redis.Set("TOKEN:"+user.Email,token,conf.ServerConf.JWTExpire * hourSecs)
	return token ,nil
}

//销毁token
func (as *AuthService)DestroyToken(email string)(bool,error){
	return redis.Del("TOKEN:"+email)
}
