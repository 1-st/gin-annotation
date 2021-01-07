package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


/* PingPong controller
[
	method:GET,
	path:/ping
]
*/
func PingPong(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"msg":"pong",
	})
}
