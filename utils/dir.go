/*
@Time : 2023/12/22 20:43
@Author : chiqing_85
@Software: GoLand
*/
package utils

import (
	"os"
	"time"
)

func Mk() string {
	day := time.Now().Format("200601")
	wd, _ := os.Getwd()
	dir := "/static/upload/" + day + "/"
	if err := os.MkdirAll(wd+dir, 0755); err != nil {
		println(err.Error())
	}
	return dir
}
