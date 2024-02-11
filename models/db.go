/*
@Time : 2023/12/16 15:25
@Author : chiqing_85
@Software: GoLand
*/
package models

import (
	"api/global"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var db *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		global.App.App.Database.Name,
		global.App.App.Database.Pas,
		global.App.App.Database.Host,
		global.App.App.Database.Port,
		global.App.App.Database.Database,
		global.App.App.Database.Charset,
	)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   global.App.App.Database.Prefix,
			SingularTable: true,
		},
	})
	if err != nil {
		fmt.Println("接连数据库失败：", err)
	}
	sqlDB, _ := db.DB()
	//SetMaxIdleConns 设置连接池中最大的闲置连接数
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置数据库的最大连接数
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置连接的最大可复用时间
	sqlDB.SetConnMaxLifetime(time.Second * 10)
}
