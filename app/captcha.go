/*
@Time : 2024/2/8 19:39
@Author : chiqing_85
@Software: GoLand
*/
package app

import (
	"api/global"
	"api/utils"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rwxrob/uniq"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var oin, _ = os.ReadDir("static/templates/original/")
var gen, _ = os.ReadDir("static/templates/gen/")

type Move struct {
	Outime float32
	X      int
	Y      []int
	Uuid   string
}

func Captcha(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())
	randhs := rand.Intn(len(gen)) + 1
	randomNumber := rand.Intn(len(oin)) + 1 // 最大值	- 目前背景图片有26个
	bg, _ := os.Open("static/templates/original/bg" + fmt.Sprintf("%d", randomNumber) + ".png")
	defer bg.Close()
	bgImg, _ := png.Decode(bg)
	hole, _ := os.Open("static/templates/gen/" + fmt.Sprintf("%d", randhs) + "/hole.png")
	defer hole.Close()
	holeImg, _ := png.Decode(hole)
	slider, _ := os.Open("static/templates/gen/" + fmt.Sprintf("%d", randhs) + "/slider.png")
	defer slider.Close()
	sliderImg, _ := png.Decode(slider)
	//bgImg
	bx := bgImg.Bounds().Dx()
	sx := sliderImg.Bounds().Dx()
	by := bgImg.Bounds().Dy()
	sy := sliderImg.Bounds().Dy()

	// 定位区域 - 最大不能超过
	w := sx/2 + rand.Intn(bx-sx-sx)
	miny := rand.Intn(by - sy)

	// 创建一个新的白色图片作为背景
	BGIMG := image.NewRGBA(bgImg.Bounds())
	draw.Draw(BGIMG, BGIMG.Bounds(), bgImg, image.ZP, draw.Src)
	draw.Draw(BGIMG, holeImg.Bounds().Add(image.Point{X: w, Y: miny}), holeImg, image.ZP, draw.Over)
	// base64编码图片 背景图
	// fmt.Println(baseImage( BGIMG ) )
	// 创建一个新的RGBA图像，用于保存抠图后的结果
	croppedImg := image.NewRGBA(holeImg.Bounds())
	// 应用遮罩层到背景图片上
	for y := 0; y < bgImg.Bounds().Dy(); y++ {
		for x := 0; x < bgImg.Bounds().Dx(); x++ {
			_, _, _, a := holeImg.At(x, y).RGBA()
			if a > 0 {
				croppedImg.Set(x, y, bgImg.At(x+w, y+miny))
			} else {
				// 如果像素完全透明，则将其设置为透明
				croppedImg.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 0})
			}
		}
	}
	draw.Draw(croppedImg, sliderImg.Bounds().Add(image.Point{X: 0, Y: 0}), sliderImg, image.ZP, draw.Over)
	uuid := uniq.Hex(10)
	global.C.Set(uuid, w, 2*time.Minute)
	// uuid
	c.JSON(http.StatusOK, gin.H{
		"uuid":       uuid,
		"top":        miny,
		"baseImage":  baseImage(BGIMG),
		"croppedImg": baseImage(croppedImg),
	})
}

func baseImage(mergedImg *image.RGBA) string {
	var buffer bytes.Buffer
	err := png.Encode(&buffer, mergedImg)
	if err != nil {
		fmt.Println("Error encoding merged image to PNG:", err)
		return ""
	}
	// 将PNG数据转换为Base64编码的字符串
	mergedBase64Str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return mergedBase64Str

}

func CaptCheck(c *gin.Context) {
	data := c.PostForm("data")
	key := c.PostForm("pubkey")
	obj := utils.AesDecrypt(data, []byte(key))
	var move Move
	json.Unmarshal(obj, &move)
	if move.Outime >= 0.2 && len(move.Y) >= 2 {
		x, found := global.C.Get(move.Uuid)
		if found && (move.X < x.(int)+5 && move.X > x.(int)-5) {
			c.JSON(http.StatusOK, gin.H{
				"message": "验证成功…",
				"result":  0,
			})
			return
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "验证失败…",
		"result":  1,
	})

}
