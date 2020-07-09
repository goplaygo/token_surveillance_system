package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
	"token_surveillance_system/extend/conf"
)

var DB *gorm.DB

func Setup(){
	var err error
	var connectString = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local",
		conf.DBConf.User,
		conf.DBConf.Password,
		conf.DBConf.Host+":"+strconv.Itoa(conf.DBConf.Port),
		conf.DBConf.DBName,
		)
	DB,err:=gorm.Open(conf.DBConf.DBType,connectString)
	if err!=nil{
		fmt.Println("mysql connect error %v",err)
		time.Sleep(10 * time.Second)//若连接失败，则延时10秒重连
		DB,err=gorm.Open(conf.DBConf.DBType,connectString)
		if err != nil {
			panic(err.Error())
		}
	}

	if DB.Error != nil {
		fmt.Println("database error %v",DB.Error)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return conf.DBConf.TablePrefix + defaultTableName
	}
	DB.LogMode(conf.DBConf.Debug)
	DB.SingularTable(true)
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
	//migrate迁移
	DB.Set(
		"gorm:table_options",
		"ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci").AutoMigrate(&User{},&Task{})
	DB.Model(&User{}).AddUniqueIndex("uk_email","email")
}
