/*
@Time : 2023/12/18 10:13
@Author : chiqing_85
@Software: GoLand
*/
package middleware

import (
	"api/global"
	"api/models"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type MyClains struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	Portrait string `json:"portrait"`
	Sex      uint8  `json:"sex"`
	City     string `json:"city"`
	Desc     string `json:"desc"`
	RegTime  int    `json:"reg_time"`
	jwt.StandardClaims
}

// 生成token
func Tokenset(u *models.User, city string) (string, error) {
	// uuid := uniq.UUID()
	var APP_KEY = []byte(global.App.App.Token.Key)
	SetClains := MyClains{
		u.Id,
		u.Nickname,
		u.Portrait,
		u.Sex,
		city,
		u.Desc,
		u.RegTime,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Hour).Unix(),
			Issuer:    global.App.App.Web.Name,
		},
	}
	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClains)
	token, err := reqClaim.SignedString(APP_KEY)
	if err != nil {
		return "", err
	}
	return token, nil
}

// 验证token
func CheckToken(token string) (*MyClains, bool) {
	var APP_KEY = []byte(global.App.App.Token.Key)
	obj, _ := jwt.ParseWithClaims(token, &MyClains{}, func(token *jwt.Token) (interface{}, error) {
		return APP_KEY, nil
	})
	if key, _ := obj.Claims.(*MyClains); obj.Valid {
		return key, true
	} else {
		return nil, false
	}
}

// token 中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHerder := c.Request.Header.Get("authorization")
		if tokenHerder == "" {
			err := errors.New("非法访问…")
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": err.Error(),
			})
			c.Abort()
			return
		}
		tokenSlice := strings.SplitN(tokenHerder, "", 2)
		if len(tokenSlice) != 2 && tokenSlice[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "非法令牌…",
			})
			c.Abort()
			return
		}
		key, ok := CheckToken(tokenHerder)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "禁止访问…",
			})
			c.Abort()
			return
		}
		if time.Now().Unix() > key.ExpiresAt {
			fmt.Println("过期了")
			c.Abort()
			return
		}
		c.Set("uid", key.Id)
		c.Next()
	}
}
