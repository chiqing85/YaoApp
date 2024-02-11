/*
@Time : 2023/12/16 0:50
@Author : chiqing_85
@Software: GoLand
*/
package global

import (
	"github.com/akyoto/cache"
	"time"
)

var C *cache.Cache

func init() {
	C = cache.New(5 * time.Minute)
	defer C.Close()
}
