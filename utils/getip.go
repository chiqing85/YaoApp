/*
@Time : 2023/12/19 21:53
@Author : chiqing_85
@Software: GoLand
*/
package utils

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"os"
)

func GetIp() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Exit(1)
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println("Local IP Address: ", ipnet.IP.String())
			}
		}
	}
}

func Geoip(ips string) string {
	wd, _ := os.Getwd()
	db, err := geoip2.Open(wd + "/GeoLite2-City.mmdb")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
	ip := net.ParseIP(string(ips))
	record, err := db.City(ip)
	if err != nil || record == nil {
		log.Fatal(err)
	}
	cou := record.Country.Names["zh-CN"]
	sub := ""
	if len(record.Subdivisions) > 0 {
		sub = record.Subdivisions[0].Names["zh-CN"]
	}
	cit := record.City.Names["zh-CN"]
	city := fmt.Sprintf("%s %s %s", cou, sub, cit)
	return city
}
