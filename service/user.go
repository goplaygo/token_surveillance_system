package service

import (
	"github.com/rs/zerolog/log"
	"token_surveillance_system/extend/code"
	"token_surveillance_system/extend/utils"
	"token_surveillance_system/models"
)

type UserService struct {
	UserID uint
	Email string
	Name string
	Password string
}

func (us *UserService)StroeUser(email string,pass string)(userID uint,err error)  {
	log.Info().Msg("storeUser Service")

	user:=&models.User{
		Email: email,
		UserName: email,
		Password: pass,
		Status: "ENABLE",
	}
	user.Password = utils.MakeSha1(user.Email + user.Password)
	log.Debug().Msgf("user password:%s",user.Password)
	userID,err = user.Insert()
	return
}

func (us *UserService)QueryByEmail(email string)(user *models.User,err error){
	userModel := &models.User{}
	condition:=map[string]interface{}{
		email: email,
	}
	user,err=userModel.FindOne(condition)
	return
}

func (us *UserService)UpdatePass(oldPass string,newPass string)(*models.User,*code.Code)  {
	userModel:=&models.User{}
	user,err:=userModel.FindOne(map[string]interface{}{"email":us.Email})
	if err != nil {
		log.Error().Msg(err.Error())
		return nil,code.ServiceInsideError
	}
	oldPassHash := utils.MakeSha1(us.Email+oldPass)
	if oldPassHash != user.Password {
		 return nil,code.AccountPassUnmatch
	}
	updateUserPass,err:=userModel.Update(user.ID, map[string]interface{}{
		"password":utils.MakeSha1(us.Email+newPass),
	})
	if err != nil {
		log.Error().Msg(err.Error())
		return nil,code.ServiceInsideError
	}
	return updateUserPass,nil
}

