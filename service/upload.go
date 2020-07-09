/*
@Time : 2020/7/9 10:50 上午
@Author : L
@File : upload.go
@Software: GoLand
*/
package service

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"token_surveillance_system/extend/conf"
	"token_surveillance_system/extend/utils"
)

type UploadService struct {

}

//GetImagePath 获取图片相对目录
func (us *UploadService)GetImagePath()string {
	return conf.ServerConf.StaticRootPath
}

//GetImgFullPath 获取图片完整目录
func (us *UploadService)GetImgFullPath()string{
	return conf.ServerConf.StaticRootPath + conf.ServerConf.UploadImagePath
}
//GetImgName 获取图片的的名称
func (us *UploadService)GetImgName(name string)string{
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name,ext)
	fileName = utils.MakeSha1(fileName)
	return fileName + ext
}

// GetImgFullURL 获取完成图片url
func (us *UploadService)GetImgFullURL(name string) string{
	return conf.ServerConf.PrefixURL + conf.ServerConf.UploadImagePath + name
}
//CheckImgExt 检查图片后缀是否满足需求
func (us *UploadService)CheckImgExt(filename string)bool {
	ext := path.Ext(filename)
	for _,allowExt := range conf.ServerConf.ImageFormats {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext){
			return true
		}
	}
	return false
}

//checkImgSize 检查图片大小是否超出
func(us *UploadService)checkImgSize(f multipart.File) bool {
	content,err := ioutil.ReadAll(f)
	if err != nil {
		log.Error().Msg(err.Error())
		return false
	}
	//单位转换
	const converRatio float64 = 1024 * 1024
	fileSize:=float64(len(content)) / converRatio
	return fileSize <= conf.ServerConf.UploadLimit
}
//checkImgPath 检查图片路径创建以及权限
func (us *UploadService)checkImgPath(path string)error {
	dir,err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.getwd err:%v",err)
	}
	isExist,err := utils.IsExist(dir + "/" +path)
	if err != nil {
		return fmt.Errorf("isExist err:%v",err)
	}
	if isExist == false {
		//若路径不存在，则创建
		err := os.MkdirAll(dir + "/" + path,os.ModePerm)
		if err != nil {
			return fmt.Errorf("mkdir err:%v",err)
		}
	}
	isPerm := utils.IsPerm(path)
	if isPerm{
		return fmt.Errorf("isPerm Permission denied src:%s",path)
	}
	return nil
}
