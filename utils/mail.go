/*
@Time : 2023/12/13 22:37
@Author : chiqing_85
@Software: GoLand
*/
package utils

import (
	"api/global"
	"gopkg.in/gomail.v2"
)

func SendMail(to string, html string, title string) error {
	host := global.App.App.Meial.Host
	port := global.App.App.Meial.Port
	user := global.App.App.Meial.From // 邮件发送者
	pwd := global.App.App.Meial.Key   // 授权码
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(user, "呦呦")) // 发送邮箱
	m.SetHeader("To", to)                            // 主送
	m.SetHeader("Subject", title)                    // 标题
	m.SetBody("text/html", html)                     // 发送html格式内容
	d := gomail.NewDialer(host, port, user, pwd)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
