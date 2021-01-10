package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/* Hello a simple controller
[
	method:GET,
	groups:/api /v1,
	path:/hello-world,
	need:auth
]
*/
func HelloWorld(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"msg": "hello, world",
	})
}

