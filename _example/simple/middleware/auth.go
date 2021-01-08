package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

/* Auth a simple middleware
[
	id:auth
]
*/
func Auth(ctx *gin.Context) {
	ctx.Next()
}

/* Log the first middleware in group /api
[
	id:log,
	group:/api@1
]
*/
func Log(ctx *gin.Context){
	fmt.Println(ctx.ClientIP())
}