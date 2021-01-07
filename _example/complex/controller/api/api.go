package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/* ApiV1 controller
[
	method:GET,
	path:/api/v1
]
*/
func ApiV1(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"pingpong":   "/ping",
		"helloworld": "/helloworld",
	})
}


/* ApiV2 controller
[
	method:GET,
	path:/api/v2,
	need:save-request-ip save-request-time
]
*/
func ApiV2(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"pingpong":   "/ping",
		"helloworld": "/helloworld",
	})
}
