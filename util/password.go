package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// 加密字符串
func EncodePassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	return string(hash)
}

// 解密字符串
func ValidatePassword(encodePassword, inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(encodePassword), []byte(inputPassword))
}
