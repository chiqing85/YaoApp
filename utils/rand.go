/*
@Time : 2023/12/15 15:10
@Author : chiqing_85
@Software: GoLand
*/
package utils

import "math/rand"

func RandomNunber() int {
	return rand.Intn(9000) + 1000
}
