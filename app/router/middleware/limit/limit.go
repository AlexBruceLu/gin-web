package limit

import (
	"fmt"
	"gin-web/app/utils/response"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func SetUp(maxBurstSize int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(time.Second), maxBurstSize)

	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
			return
		}
		fmt.Println("too many request")

		utilGin := response.Gin{Ctx: c}
		utilGin.Response(-1, "too-many-request", nil)
		c.Abort()
		return

	}

}
