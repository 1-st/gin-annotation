package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

/* SaveRequestIP _example middleware
[
	id:save-request-ip
]
*/
func SaveRequestIP(ctx *gin.Context) {
	fmt.Println(ctx.ClientIP())
	ctx.Next()
}

/*	SaveRequestTime _example middleware
[
	id:save-request-time
]
*/
func SaveRequestTime(ctx *gin.Context) {
	fmt.Println(time.Now().String())
	ctx.Next()
}
