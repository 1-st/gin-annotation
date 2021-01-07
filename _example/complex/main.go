package main

import (
	"github.com/gin-gonic/gin"
)


func main() {
	e := gin.Default()
	Route(e)
	_ = e.Run("0.0.0.0:8080")
}
