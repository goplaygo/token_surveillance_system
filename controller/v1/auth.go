package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"token_surveillance_system/extend/code"
	"token_surveillance_system/extend/jwt"
	"token_surveillance_system/extend/utils"
	"token_surveillance_system/service"
)

//用户控制器
type AuthController struct{}

//SignupRequest 账号注册请求参数
type SignupRequest struct {
	Email string `json:"email" binding:"required,email"`
	AccountPass string `json:"accountPass" binding:"required"`
	ConfirmPass string `json:"confirmPass" binding:"required"`
}

func(au AuthController)Signup(c *gin.Context){
	log.Info().Msg("signup controller")
	reqBody := SignupRequest{}
	if err :=c.ShouldBindJSON(&reqBody);err!=nil {
		log.Error().Msg(err.Error())
		utils.ResponseFormat(c,code.RequestParamError,nil)
		return
	}
	log.Debug().Msgf("email param:%s",reqBody.Email)
	log.Debug().Msgf("confirmPass param:%s",reqBody.ConfirmPass)

	if reqBody.AccountPass != reqBody.ConfirmPass {
		utils.ResponseFormat(c,code.SignupPassUnmatch,nil)
		return
	}

	userService := service.UserService{
		Email: reqBody.Email,
		Password: reqBody.ConfirmPass,
	}

	userID,err := userService.StroeUser(reqBody.Email,reqBody.ConfirmPass)
	if err != nil {
		log.Error().Msg(err.Error())
		utils.ResponseFormat(c,code.ServiceInsideError,nil)
		return
	}
	log.Info().Msgf("signup controller result userId:%s",userID)
	utils.ResponseFormat(c,code.Success,map[string]uint{"userID":userID})
	return

}

type SigninRequest struct {
	Email string `json:"email"binding:"required,email"`
	Password string `json:"password"binding:"required,max=20"`
}

func (ac AuthController)Signin(c *gin.Context) {
	log.Info().Msg("signin controller")
	reqBody := SigninRequest{}
	err := c.ShouldBindJSON(&reqBody)
	if err != nil {
		log.Error().Msg(err.Error())
		utils.ResponseFormat(c, code.SigninInfoError, nil)
		return
	}
	//登录验证
	userService := service.UserService{
		Email: reqBody.Email,
	}
	user, err := userService.QueryByEmail(reqBody.Email)
	if err != nil {
		log.Error().Msg(err.Error())
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}

	if user == nil || user.Password != utils.MakeSha1(reqBody.Email + reqBody.Password){
		utils.ResponseFormat(c,code.SigninInfoError,nil)
		return
	}

	//生成token
	authService := service.AuthService{
		User: user,
	}
	token,err:=authService.GenerateToken(*user)
	if err!=nil{
		utils.ResponseFormat(c,code.ServiceInsideError,nil)
		return
	}
	utils.ResponseFormat(c,code.Success,map[string]interface{}{
		"userId":user.ID,
		"userName":user.UserName,
		"email":user.Email,
		"token":token,
	})
	return
}

func (ac AuthController)Signout(c *gin.Context)  {
	log.Info().Msg("signout controller")
	claims := c.MustGet("claims").(*jwt.CustomClaims)
	log.Debug().Msgf("claims:%v",claims)
	//销毁token
	authService := service.AuthService{}
	isOK,err:=authService.DestroyToken(claims.Email)
	if err!=nil || isOK == false{
		utils.ResponseFormat(c,code.ServiceInsideError,nil)
		return
	}
	utils.ResponseFormat(c,code.Success,map[string]interface{}{})
	return
}
