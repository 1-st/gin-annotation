package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/* Hello hello world controller
[
	method:Any,
	path:/hello-world,
	groups:/extra
]
*/
func HelloWorld(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"msg": "hello, world",
	})
}

/* Hello hello world controller
[
	method:GET,
	path:/hello-world
]
*/
func HelloWorldFake1(a string) {

}

/* Hello hello world controller
[
	method:GET,
	path:/hello-world
]
*/
func HelloWorldFake2() bool {
	return true
}
