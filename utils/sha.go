/*
@Time : 2023/12/17 0:16
@Author : chiqing_85
@Software: GoLand
*/
package utils

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// GeneratePassword 给密码就行加密操作
func Generate(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

// ValidatePassword 密码比对
func ValidatePassword(userPassword string, hashed string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return errors.New("密码错误…")
	}
	return nil

}
