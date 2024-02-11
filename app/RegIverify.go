/*
@Time : 2023/12/15 14:00
@Author : chiqing_85
@Software: GoLand
*/
package app

import (
	"api/global"
	"api/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func RegSnd(c *gin.Context) {
	value := c.PostForm("value")
	reg := "^1[345789]\\d{9}"
	rand := utils.RandomNunber()
	if regexp.MustCompile(reg).MatchString(value) {
		// 短信模板
		fmt.Println("phone")
	} else {
		// 邮件模板
		file_tmp := "./static/tmp/regi_mail.html"
		cont, err := os.ReadFile(file_tmp)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		str := strings.Replace(string(cont), "{rand}", strconv.Itoa(rand), -1)
		str = strings.Replace(str, "{name}", global.App.App.Web.Name, -1)
		str = strings.Replace(str, "{url}", global.App.App.Web.Url, -1)
		title := "用户注册" // 邮件标题
		go utils.SendMail(value, str, title)
		global.C.Set(value, strconv.Itoa(rand), 10*time.Minute)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "邮件发送成功…",
		})
	}
}
