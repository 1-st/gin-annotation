package v2

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/* ID user private id version 2
[
	method:GET,
	group:/api/v2 /user/:name /private,
	path:/id
]
*/
func ID(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"id":      "100001",
		"version": "2",
	})
}
