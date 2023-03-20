package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Res struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	Success = 1
	Fail    = 0
)

//R.Ok(c, "自定义msg",data)
func Ok(c *gin.Context, msg string, data interface{}) {
	Response(c, http.StatusOK, Success, msg, data)
}

//R.Error(c, "自定义msg",data)
func Error(c *gin.Context, msg string, data interface{}) {
	Response(c, http.StatusOK, Fail, msg, data)
}

//R.Response(c,200,1,"msg",data)
func Response(c *gin.Context, status int, code int, msg string, data interface{}) {
	c.JSON(status, Res{
		code,
		msg,
		data,
	})
}
