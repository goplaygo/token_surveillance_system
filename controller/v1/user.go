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

type UserNameRequest struct {
	Name string `json:"name"binding:"required"`
}

func (uc UserController)AlterName(c *gin.Context)  {
	log.Info().Msg("user change name controller")
	claims :=c.MustGet("claims").(*jwt.CustomClaims)
	if claims == nil {
		utils.ResponseFormat(c,code.TokenInvalid,nil)
		return
	}
	reqBody := UserNameRequest{}
	if err:=c.ShouldBindJSON(&reqBody);err!=nil{
		utils.ResponseFormat(c,code.RequestParamError,nil)
		return
	}
	userService:=service.UserService{UserID: claims.ID}
	updateUser,msgCode:=userService.UpdateName(reqBody.Name)
	if msgCode != nil {
		utils.ResponseFormat(c,msgCode,nil)
		return
	}
	utils.ResponseFormat(c,code.Success,map[string]interface{}{
		"userId":updateUser.ID,
		"userName":updateUser.UserName,
	})
}

type UserPassRequest struct {
	OldPass string `json:"oldPass"binding:"required,max=20"`
	NewPass string `json:"newPass"binding:"required,max=20"`
}

func (uc UserController)AlterPass(c *gin.Context){
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

func (uc UserController)AlterAvatar(c *gin.Context){
	log.Info().Msg("user change avatar")
	claims:=c.MustGet("claims").(*jwt.CustomClaims)
	if claims == nil {
		utils.ResponseFormat(c,code.TokenInvalid,nil)
		return
	}

	//获取文件上传内容
	file,image,err := c.Request.FormFile("avatar")
	if err!=nil{
		log.Error().Msg(err.Error())
		utils.ResponseFormat(c,code.ServiceInsideError,nil)
		return
	}
	if image == nil {
		utils.ResponseFormat(c,code.RequestParamError,nil)
		return
	}
	//获取头像名称
	uploadService := service.UploadService{}
	avatarName := uploadService.GetImgName(image.Filename)
	fullPath := uploadService.GetImgFullPath()
	//检验图片格式
	if !uploadService.CheckImgExt(avatarName){
		utils.ResponseFormat(c,code.UploadSuffixError,nil)
		return
	}
	//检查图片大小
	if !uploadService.CheckImgSize(file){
		utils.ResponseFormat(c,code.UploadSizeLimit,nil)
		return
	}
	err = uploadService.CheckImgPath(fullPath)
	if err != nil {
		log.Error().Msg(err.Error())
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}
	err = c.SaveUploadedFile(image, fullPath+avatarName)
	if err != nil {
		log.Error().Msg(err.Error())
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}

	//更新图片
	userService:=service.UserService{UserID: claims.ID}
	updateUserAvatar,msgCode:=userService.UpdateAvatar(uploadService.GetImagePath()+avatarName)

	if msgCode!=nil{
		utils.ResponseFormat(c,msgCode,nil)
		return
	}
	utils.ResponseFormat(c,code.Success,map[string]interface{}{
		"userId":updateUserAvatar.ID,
		"avatar":updateUserAvatar.Avatar,
	})
}
