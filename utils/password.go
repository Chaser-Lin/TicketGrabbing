package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// 使用 bcrypt包的加密算法对密码进行加密
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密码加密失败:%v", err)
	}
	return string(hashedPassword), nil
}

// 验证 hashPassword 是不是由 password 加密生成的
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
