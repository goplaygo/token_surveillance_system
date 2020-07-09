package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"token_surveillance_system/extend/code"
	"token_surveillance_system/extend/jwt"
	"token_surveillance_system/extend/utils"
	"token_surveillance_system/service"
)

type UserController struct {

}


func (uc UserController)FindUser(c *gin.Context){
	claims := c.MustGet("claims").(*jwt.CustomClaims)
	if claims != nil {
		utils.ResponseFormat(c,code.Success,map[string]interface{}{
			"data":claims,
		})
		return
	}
}

type UserPassRequest struct {
	OldPass string `json:"oldPass"binding:"required,max=20"`
	NewPass string `json:"newPass"binding:"required,max=20"`
}

func (uc UserController)Alterpass(c *gin.Context){
	log.Info().Msg("change password controller")
	claims:=c.MustGet("claims").(*jwt.CustomClaims)
	if claims == nil {
		utils.ResponseFormat(c,code.TokenInvalid,nil)
		return
	}
	//获取请求参数
	reqBody:=UserPassRequest{}
	if err:=c.ShouldBindJSON(&reqBody);err!=nil{
		utils.ResponseFormat(c,code.RequestParamError,nil)
		return
	}
	userService := service.UserService{Email: claims.Email}
	updateUser,msgCode :=userService.UpdatePass(reqBody.OldPass,reqBody.NewPass)
	if msgCode!=nil{
		utils.ResponseFormat(c,msgCode,nil)
		return
	}
	utils.ResponseFormat(c,code.Success,map[string]interface{}{
		"userId":updateUser.ID,
		"userName":updateUser.UserName,
		"email":updateUser.Email,
	})
}
