/*
@Time : 2020/7/7 1:30 下午
@Author : L
@File : jwt.go
@Software: GoLand
*/

package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
	"token_surveillance_system/extend/conf"
)

//JWT认证相关
type JWT struct {
	JWTSecret []byte
}

//创建JWT实例
func NewJWT() *JWT  {
	return &JWT{[]byte(conf.ServerConf.JWTSecret)}
}

var (
	//ErrTokenExpired 令牌验证失败
	ErrTokenExpired = errors.New("Token is expired")
	//ErrTokenNotVaildYet 验证令牌未激活
	ErrTokenNotVaildYet = errors.New("Token not active yet")
	//ErrTokenMalformed 验证并非令牌
	ErrTokenMalformed = errors.New("That is not even a token")
	//ErrTokenInvaild 验证无效令牌
	ErrTokenInvaild = errors.New("could not handle this token")
)

type CustomClaims struct {
	Id uint `json:"userId"`
	UserName string `json:"userName"`
	Email string `json:"email"`
	jwt.StandardClaims
}

//生成token
func (j *JWT)CreatToken(claims CustomClaims) (string, error){
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodES256,claims)
	return tokenClaims.SignedString(j.JWTSecret)
}

//解析token
func (j *JWT)ParseToken(token string) (*CustomClaims,error) {
	tokenClaims ,err := jwt.ParseWithClaims(token,&CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.JWTSecret,nil
	})
	if err !=nil{
		if ve,ok:=err.(jwt.ValidationError);ok{
			if ve.Errors&jwt.ValidationErrorMalformed!=0{
				return nil,ErrTokenMalformed
			}else if ve.Errors&jwt.ValidationErrorExpired!=0{
				return nil,ErrTokenExpired
			}else if ve.Errors&jwt.ValidationErrorNotValidYet!=0{
				return nil,ErrTokenNotVaildYet
			} else {
				return nil,ErrTokenInvaild
			}
		}
	}
	if clamis,ok:=tokenClaims.Claims.(*CustomClaims);ok&&tokenClaims.Valid{
		return clamis,nil
	}
	return nil,ErrTokenInvaild
}


//刷新token
func (j *JWT)RefreshToken(token string) (string,error)  {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0,0)
	}
	tokenClaims,err:=jwt.ParseWithClaims(token,&CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.JWTSecret,nil
	})
	if err!=nil{
		return "", err
	}
	if claims,ok:=tokenClaims.Claims.(*CustomClaims);ok && tokenClaims.Valid{
		jwt.TimeFunc = time.Now
		expiredTime := time.Duration(conf.ServerConf.JWTExpire)
		claims.StandardClaims.ExpiresAt = time.Now().Add(expiredTime * time.Hour).Unix()
		return j.CreatToken(*claims)
	}
	return "",ErrTokenInvaild
}
