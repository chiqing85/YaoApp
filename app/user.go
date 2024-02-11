/*
@Time : 2023/12/16 11:37
@Author : chiqing_85
@Software: GoLand
*/
package app

import (
	"api/global"
	"api/middleware"
	"api/models"
	"api/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 用户登录
func Login(c *gin.Context) {
	name := c.PostForm("name")
	pas := c.PostForm("pas")
	if name == "" || pas == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"code": -1,
			"msg":  "禁止非法访问…",
		})
		return
	}
	u, e := models.Login(name)
	if e != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": -1,
			"msg":  e.Error(),
		})
		return
	}
	if err := utils.ValidatePassword(pas, u.Pas); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	clientIP := c.ClientIP()
	if clientIP == "127.0.0.1" {
		clientIP = "223.104.61.37"
	}
	city := utils.Geoip(clientIP)
	token, err := middleware.Tokenset(u, city)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "恭喜，登录成功",
		"token": token,
	})
}

// 注册用户
func Regi(c *gin.Context) {
	mail := c.PostForm("mail")
	phone := c.PostForm("phone")
	condition := ""
	key := ""
	if len(mail) > 0 {
		condition = "mail"
		key = mail
	}
	if len(phone) > 0 {
		condition = "phone"
		key = phone
	}
	verify := c.PostForm("verify")
	pas := c.PostForm("pas")
	if len(verify) > 0 && len(pas) > 0 {
		v, _ := global.C.Get(key)
		if v == nil || v != verify {
			c.JSON(http.StatusNotFound, gin.H{
				"code": -1,
				"msg":  "验证码错误…",
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  "非法访问…",
		})
		return
	}
	found := models.CheckUser(condition, key)
	if found {
		var t = ""
		if condition == "mail" {
			t = "邮箱"
		} else {
			t = "手机"
		}
		c.JSON(http.StatusNotFound, gin.H{
			"code": -1,
			"msg":  fmt.Sprintf("该%s号已被注册…", t),
		})
		return
	}
	pass, _ := utils.Generate(pas)
	e := models.UserReg(condition, key, string(pass))
	if e != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": -1,
			"msg":  e.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "恭喜，用户注册成功,可以去登录了…",
		})
	}
}

func CartUpdate(c *gin.Context) { // 用户卡片资料修改
	data := c.PostFormMap("data")
	uid, _ := c.Get("uid")
	if len(data) > 0 {
		u, err := models.CartUpdate(data, uid)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code": -1,
				"msg":  err.Error(),
			})
			return
		}
		clientIP := c.ClientIP()
		if clientIP == "127.0.0.1" {
			clientIP = "223.104.61.37"
		}
		city := utils.Geoip(clientIP)
		token, e := middleware.Tokenset(u, city)
		if e != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code": -1,
				"msg":  e.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  200,
			"msg":   "更新成功…",
			"token": token,
		})
	}
}
