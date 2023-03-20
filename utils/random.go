package utils

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var alpha = "QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm"
var digit = "0123456789"

func RandomString(min int) string {
	len := rand.Intn(5) + min
	return RandomAlpha(len)
}

func RandomAlpha(len int) string {
	sb := strings.Builder{}
	for ; len > 0; len-- {
		sb.WriteByte(alpha[rand.Intn(52)])
	}
	return sb.String()
}

func RandomEmail() string {
	sb := strings.Builder{}
	len := 10
	for ; len > 0; len-- {
		sb.WriteByte(digit[rand.Intn(10)])
	}
	sb.WriteString("@qq.com")
	return sb.String()
}
