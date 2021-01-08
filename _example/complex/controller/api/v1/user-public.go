package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/* Age user age controller
[
	method:GET,
	groups:/api/v1 /user/:name,
	path:/age
]
*/
func Age(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]int{
		"age":     18,
		"version": 1,
	})
}

/* Avatar user Avatar controller
[
	method:GET,
	groups:/api/v1 /user/:name,
	path:/avatar
]
*/
func Avatar(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"avatar":  "_example.png",
		"version": "1",
	})
}
