package middleware

import (
	"Project/MyProject/response"
	"Project/MyProject/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	AdminID   = 1
	AdminName = "admin"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Res{
				Code: response.Fail,
				Msg:  "缺少authorization头部字段",
			})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Res{
				Code: response.Fail,
				Msg:  "authorization头部字段格式不正确",
			})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationType {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Res{
				Code: response.Fail,
				Msg:  fmt.Sprintf("不支持的authorization头部类型：%s", authorizationType),
			})
			return
		}

		accessToken := fields[1]
		accessPayload, err := utils.TokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Res{
				Code: response.Fail,
				Msg:  err.Error(),
			})
			return
		}

		if accessPayload.UserID != AdminID || accessPayload.Email != AdminName {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Res{
				Code: response.Fail,
				Msg:  fmt.Sprintf("登录的用户(%s)并非管理员，无法进行此操作!", accessPayload.Email),
			})
			return
		}

		c.Set(TokenKey, accessToken)
		c.Set(AuthorizationPayloadKey, accessPayload)
		c.Set(UserID, accessPayload.UserID)
		c.Next()
	}
}
