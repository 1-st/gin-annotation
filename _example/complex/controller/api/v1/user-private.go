package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


/* ID user private id version 1
[
	method:GET,
	groups:/api/v1 /user/:name /private,
	path:/id
]
*/
func ID(ctx *gin.Context){
	ctx.JSON(http.StatusOK,map[string]string{
		"id":"100001",
		"version":"1",
	})
}
