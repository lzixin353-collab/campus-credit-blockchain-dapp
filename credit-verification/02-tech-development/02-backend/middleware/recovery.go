// middleware/recovery.go 自定义 Recovery，panic 时返回 JSON 便于前端显示具体错误
package middleware

import (
	"fmt"
	"log"
	"net/http"

	"campus-credit-backend/utils"

	"github.com/gin-gonic/gin"
)

// RecoveryJSON 捕获 panic 并返回统一 JSON，方便前端展示 msg
func RecoveryJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Recovery] panic: %v", r)
				// 若 handler 已写过响应则不再写，避免二次写导致 500
				if !c.Writer.Written() {
					c.AbortWithStatusJSON(http.StatusOK, utils.Response{
						Code: 500,
						Msg:  "服务器异常: " + toString(r),
						Data: nil,
					})
				}
			}
		}()
		c.Next()
	}
}

func toString(v interface{}) string {
	switch x := v.(type) {
	case string:
		return x
	case error:
		return x.Error()
	default:
		return fmtString(v)
	}
}

func fmtString(v interface{}) string {
	defer func() { _ = recover() }()
	return fmt.Sprintf("%v", v)
}
