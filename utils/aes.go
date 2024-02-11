/*
@Time : 2024/2/10 21:13
@Author : chiqing_85
@Software: GoLand
*/
package utils

import (
	"crypto/aes"
	"encoding/base64"
)

// 加密
func Aesecrypt() {

}

// 解密
func AesDecrypt(cypted string, key []byte) []byte {
	data, _ := base64.StdEncoding.DecodeString(cypted)
	block, _ := aes.NewCipher(key)
	decrypted := make([]byte, len(data))
	size := block.BlockSize()
	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Decrypt(decrypted[bs:be], data[bs:be])
	}
	return PKCS5UnPadding(decrypted)
}

func PKCS5UnPadding(decrypted []byte) []byte {
	length := len(decrypted)
	if length == 0 {
		return decrypted
	}
	// 2. 根据填充的字节值得到填充的次数
	number := int(decrypted[length-1])

	return decrypted[:(length - number)]
}
