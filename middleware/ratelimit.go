package middleware

import (
	"Project/MyProject/response"
	"github.com/gin-gonic/gin"

	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

func RateLimit(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			response.Response(c, http.StatusForbidden, 0, "短时间内请求太多，当前接口已限流", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
