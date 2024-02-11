/*
@Time : 2023/12/15 14:24
@Author : chiqing_85
@Software: GoLand
*/
package router

import (
	"api/app"
	"api/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

//type RouterGroup struct {
//	prefix string
//	engine *Engine
//	parent *RouterGroup
//}
//type Engine struct {
//	*RouterGroup
//}
//
//func (e *Engine) Run(port ...string) {
//	var Port string
//
//	err := http.ListenAndServe(Port, nil)
//	if err != nil {
//		fmt.Printf("http server failed, err:%v\n", err)
//	}
//
//}
//
//func cors(f http.HandlerFunc) http.HandlerFunc {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Access-Control-Allow-Origin", "*")                                                                                         // 指明哪些请求源被允许访问资源，值可以为 "*"，"null"，或者单个源地址。
//		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")                              //对于预请求来说，指明了哪些头信息可以用于实际的请求中。
//		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")                                                                       //对于预请求来说，哪些请求方式可以用于实际的请求。
//		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type") //对于预请求来说，指明哪些头信息可以安全的暴露给 CORS API 规范的 API
//		w.Header().Set("Access-Control-Allow-Credentials", "true")                                                                                 //指明当请求中省略 creadentials 标识时响应是否暴露。对于预请求来说，它表明实际的请求中可以包含用户凭证。
//
//		//放行所有OPTIONS方法
//		if r.Method == "OPTIONS" {
//			w.WriteHeader(200)
//			return
//		}
//		f.ServeHTTP(w, r)
//	})
//}
//
//func New() *Engine {
//	engine := &Engine{}
//	engine.RouterGroup = &RouterGroup{}
//	return engine
//
//	//http.HandleFunc("/verify", cors(verify))
//	//err := http.ListenAndServe("127.0.0.1:80", nil)
//	//if err != nil {
//	//	fmt.Printf("http server failed, err:%v\n", err)
//	//	return nil
//	//}
//	//return nil
//}

func Cors() gin.HandlerFunc {
	//return func(ctx *gin.Context) {
	//	//以下是引用跨域 cookie
	//	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	//	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,Content-Type")
	//	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
	//	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	//}
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" { //放行所有OPTIONS方法
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func Router(r *gin.Engine) *gin.Engine {
	r.Use(Cors())
	r.POST("/newchat", app.Chat)
	r.POST("/verify", app.RegSnd)
	r.POST("/reg", app.Regi)
	r.POST("/login", app.Login)
	r.POST("/captcha", app.Captcha)
	r.POST("/captcha/check", app.CaptCheck)
	r.POST("/upload/:branch", middleware.JwtToken(), app.Upload)
	r.POST("/cartdesc", middleware.JwtToken(), app.CartUpdate)

	return r
}
