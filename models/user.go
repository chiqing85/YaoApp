/*
@Time : 2023/12/16 19:31
@Author : chiqing_85
@Software: GoLand
*/
package models

import (
	"errors"
	"gitlab.com/rwxrob/uniq"
	"regexp"
	"time"
)

type User struct {
	Id        int    `json:"id" gorm:"PRIMARY_KEY;auto_increment"`
	Nickname  string `json:"nickname" gorm:"type:varchar(60);unique_index"`
	Pas       string `json:"pas" gorm:"type:varchar(60)"`
	Portrait  string `json:"portrait" gorm:"type:varchar(255)"`
	Phone     string `json:"phone" gorm:"type:varchar(11)"`
	Mail      string `json:"mail" gorm:"type:varchar(60)"`
	Sex       uint8  `json:"sex" gorm:"type:tinyint(1);default 0;not null"`
	Desc      string `json:"desc" gorm:"type:varchar(30)"`
	RegTime   int    `json:"reg_time"  gorm:"type:int"`
	LoginTime int    `json:"login_time" gorm:"type:int"`
	LoginIp   string `json:"login_ip" gorm:"type:varchar(16)"`
	LockState uint8  `json:"lock_state" gorm:"type:tinyint(1);default 0;not null"`
}

func CheckUser(cond, key string) bool {
	var user User
	db.Select("id").Where(cond+" = ?", key).First(&user)
	if user.Id > 0 {
		return true
	}
	return false
}

func Login(n string) (*User, error) { // 登录
	var user User
	db.Where("mail = ?", n).Or("nickname = ?", n).Or("phone = ?", n).Take(&user)
	if user.Id == 0 {
		return nil, errors.New("用户名不存在…")
	}
	return &user, nil
}

// 注册
func UserReg(cond, key, pas string) error {
	t := time.Now().Unix()
	user := User{
		Pas:       pas,
		RegTime:   int(t),
		LockState: 1,
		Nickname:  uniq.Hex(4),
	}
	reg := "^1[345789]\\d{9}"
	if regexp.MustCompile(reg).MatchString(cond) {
		user.Phone = key
	} else {
		user.Mail = key
	}
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func CartUpdate(data map[string]string, uid any) (*User, error) { // 更新用户资料
	var user User
	db.First(&user, uid)
	for k, v := range data {
		if e := db.Model(&user).Update(k, v).Error; e != nil {
			return nil, e
		}
	}
	return &user, nil
}

func PortraitUpdate(p string, uid any) (*User, error) { // 更新用户头像
	var user User
	db.First(&user, uid)
	user.Portrait = p
	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
