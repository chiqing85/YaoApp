/*
@Time : 2023/12/22 18:14
@Author : chiqing_85
@Software: GoLand
*/
package app

import (
	"api/middleware"
	"api/models"
	"api/utils"
	"bytes"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var picMap map[string]string

func init() {
	picMap = make(map[string]string)
	picMap["0000010001"] = "ico"
	picMap["ffd8ffe0"] = "jpeg"
	picMap["ffd8ffe1"] = "jpg"
	picMap["ffd8ffee"] = "jpg"
	picMap["ffd8ffe8"] = "jpg"
	picMap["89504e47"] = "png"
	picMap["47494638"] = "gif"
	picMap["52494646"] = "webp"
	picMap["49492a00"] = "tif"
	picMap["424d"] = "bmp"
	picMap["41433130"] = "dwg"
	picMap["3c21444f"] = "html"
	picMap["3c68746d"] = "html"
	picMap["3c21646f"] = "html"
	picMap["48544d4c"] = "css"
	picMap["2a207b0d"] = "css"
	picMap["696b2e71"] = "js"
	picMap["636c6173"] = "js"
	picMap["7b5c7274"] = "rtf"
	picMap["38425053"] = "psd"
	picMap["46726f6d"] = "eml"
	picMap["d0cf11e0"] = "doc"
	picMap["d0cf11e0"] = "vsd"
	picMap["5374616E"] = "mdb"
	picMap["25215053"] = "ps"
	picMap["25504446"] = "pdf"
	picMap["2e524d46"] = "rmvb"
	picMap["464c5601"] = "flv"
	picMap["00000020"] = "mp4"
	picMap["49443303"] = "mp3"
	picMap["000001ba"] = "mpg"
	picMap["3026b275"] = "wmv"
	picMap["52494646e27807005741"] = "wav"
	picMap["52494646d07d60074156"] = "avi"
	picMap["4d546864"] = "mid"
	picMap["504b030414"] = "zip"
	picMap["52617221"] = "rar"
	picMap["23546869"] = "ini"
	picMap["504b03040a"] = "jar"
	picMap["4d5a9000"] = "exe"
	picMap["3c254020"] = "jsp"
	picMap["4d616e69"] = "mf"
	picMap["3c3f786d"] = "xml"
	picMap["494e5345"] = "sql"
	picMap["7061636b"] = "java"
	picMap["40656368"] = "bat"
	picMap["1f8b0800"] = "gz"
	picMap["6c6f6734"] = "properties"
	picMap["cafebabe"] = "class"
	picMap["49545346"] = "chm"
	picMap["04000000"] = "mxp"
	picMap["504b0304"] = "docx"
	picMap["d0cf11e0"] = "wps"
	picMap["6431303a"] = "torrent"
	picMap["6D6F6F76"] = "mov"
	picMap["FF575043"] = "wpd"
	picMap["CFAD12FE"] = "dbx"
	picMap["2142444E"] = "pst"
	picMap["AC9EBD8F"] = "qdf"
	picMap["E3828596"] = "pwl"
	picMap["2E7261FD"] = "ram"
}

func Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "上传文件失败…",
		})
		return
	}
	branch := c.Param("branch")
	if branch == "portrait" { // 分支
		if header.Size > 1024*100 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": -1,
				"msg":  "图片不能大于100KB…",
			})
			return
		}
		if t, ok := Portrait(file); ok {
			// 存入文件
			// 创建当月文件夹
			dir := utils.Mk()
			// 生成文件名称
			FileName := strconv.FormatInt(time.Now().UnixNano(), 10)
			save := c.SaveUploadedFile(header, "."+dir+FileName+"."+t)
			if save == nil {
				//
				uid, _ := c.Get("uid")
				u, e := models.PortraitUpdate(dir+FileName+"."+t, uid)
				if e != nil {
					c.JSON(http.StatusNotFound, gin.H{
						"code": -1,
						"msg":  e.Error(),
					})
					return
				}
				clientIP := c.ClientIP()
				if clientIP == "127.0.0.1" {
					clientIP = "223.104.61.37"
				}
				city := utils.Geoip(clientIP)
				token, er := middleware.Tokenset(u, city)
				if e != nil {
					c.JSON(http.StatusNotFound, gin.H{
						"code": -1,
						"msg":  er.Error(),
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"code":  200,
					"msg":   "图片上传成功",
					"token": token,
					"path":  dir + FileName + "." + t,
				})
				return
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"code": -1,
				"msg":  "文件非图片格式…",
			})
		}
	}

}
func Portrait(f multipart.File) (string, bool) {
	cont, err := io.ReadAll(f)
	if err != nil {
		return "", false
	}
	t, b := FileType(cont[:10])
	if b && IsImage(t) {
		return t, true
	}
	return t, false
}
func FileType(cont []byte) (string, bool) {
	if cont == nil || len(cont) <= 0 {
		return "", false
	}
	res := bytes.Buffer{}
	temp := make([]byte, 0)
	for _, v := range cont {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	var filetype string
	for s, v := range picMap {
		if strings.HasPrefix(res.String(), s) {
			filetype = v
			break
		}
	}
	return filetype, true
}
func IsImage(t string) bool {
	allow := map[string]bool{
		"jpg":  true,
		"png":  true,
		"gif":  true,
		"jpeg": true,
		"webp": true,
		"bmp":  true,
	}
	_, ok := allow[t]
	if !ok {
		return false
	}
	return true
}
