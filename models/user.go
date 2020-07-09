package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	UserName string `gorm:"column:name;type:varchar(255);unique_index;default:null"`
	Password string `gorm:"column:password;type:varchar(255);default:null"`
	Email    string `gorm:"column:email;type:varchar(255);unique_index;default:null"`
	Avatar   string `gorm:"column:avatar;type:varchar(255);default:null"`
	Status   string `sql:"type:ENUM('ENABLE','DISABLE')"`
}

//插入用户
func (user *User) Insert() (userID uint, err error) {
	result := DB.Create(&user)
	userID = user.ID
	if result.Error != nil {
		err = result.Error
	}
	return
}

func (user *User) FindOne(condition map[string]interface{}) (*User, error) {
	var userInfo User
	result := DB.Select("id,name,email,avatar,password").Where(condition).Find(&userInfo)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	if userInfo.ID > 0 {
		return &userInfo, nil
	}
	return nil, nil
}

func (user *User) FindAll(pageNum int, pageSize int, condition interface{}) (users []User, err error) {
	result := DB.Offset(pageNum).Limit(pageSize).Select("id,name,email,avatar").Where(condition).Find(&users)
	err = result.Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	return
}

func (user *User) Update(userID uint, data map[string]interface{}) (*User, error) {
	err := DB.Model(&User{}).Where("id=?", userID).Updates(data).Error
	if err != nil {
		return nil, err
	}
	var updateUser User
	err = DB.Select([]string{"id", "name", "email", "avatar"}).First(&updateUser, userID).Error
	if err != nil {
		return nil, err
	}
	return &updateUser, nil
}

func (user *User)Delete(userID uint)(delUser User,err error){
	if err = DB.Select([]string{"id"}).Find(&user,userID).Error;err!=nil{
		return
	}
	if err = DB.Delete(&user).Error;err!=nil{
		return
	}
	delUser = *user
	return
}
