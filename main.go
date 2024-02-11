/*
@Time : 2023/12/13 17:49
@Author : chiqing_85
@Software: GoLand
*/
package main

import (
	"api/global"
	"api/models"
	"api/router"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	global.InitConfig()
	models.InitDB()
	r := gin.Default()
	r.Static("static", "./static")
	e := router.Router(r)
	e.Run(":80")
}
