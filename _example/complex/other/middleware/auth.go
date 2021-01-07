package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

/* Auth1 _example middleware
[
	id:auth1,
	groups:/api/v1/user/:name/private@1,
	groups:/api/v2/user/:name/private@1
]
*/
func Auth1(ctx *gin.Context) {
	name, ok := ctx.Params.Get("name")
	if ok {
		fmt.Println(name + "auth1 success")
		ctx.Next()
	} else {
		fmt.Println("auth1 failed")
		ctx.Abort()
	}
}

/* Auth2 _example middleware
[
	id:auth2,
	groups:/api/v1/user/:name/private@2,
	groups:/api/v2/user/:name/private@2
]
*/
func Auth2(ctx *gin.Context) {
	name, ok := ctx.Params.Get("name")
	if ok {
		fmt.Println(name + "auth2 success")
		ctx.Next()
	} else {
		fmt.Println("auth2 failed")
		ctx.Abort()
	}
}

