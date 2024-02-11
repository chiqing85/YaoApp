/*
@Time : 2023/12/15 22:49
@Author : chiqing_85
@Software: GoLand
*/
package global

import (
	"api/conf"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Appconf struct {
	App *conf.Config
}

var App = new(Appconf)

func InitConfig() {
	var cfgfile = "config.yaml"
	data, err := os.ReadFile(cfgfile)
	if err != nil {
		fmt.Println(err.Error())
	}
	if e := yaml.Unmarshal(data, &App.App); err != nil {
		fmt.Println(e.Error())
	}
}
