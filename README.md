# gin-annotation

gin-annotation is a powerful cli tool to implement gin annotation

 <img src="https://raw.githubusercontent.com/1-st/logos/master/gin-annotation/logo.png" width = "50%" alt="logo" align=center /> 

## Features

* Using code generating technology by operating golang [AST](https://en.wikipedia.org/wiki/Abstract_syntax_tree)
* Routing group
* Enable Aspect Oriented Programming for gin by middleware annotation

## Quick Start

1. Installation.

```shell
go get github.com/1-st/gin-annotation
```

2. Write your HandlerFunc anywhere.

> source code files see dir: *_example/simple*

```go
// controller/hello.go
package controller

/* Hello a simple controller
[
	method:GET,
	groups:/api,
	path:/hello-world,
	need:auth
]
*/
func HelloWorld(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"msg": "hello, world",
	})
}
```

```go
// middleware/auth.go
package middleware

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
func Log(ctx *gin.Context) {
	fmt.Println(ctx.ClientIP())
}
```

3. Run gin-annotation at dir: *_example/simple*:

```sh
$ gin-annotation ./controller ./middleware
$ ls
controller main.go route.entry.go(!!!new file)
```

> tips: the name of new file is decided by environment variable <font color=#ee00ee>GIN_ANNOTATION_FILE</font>
, default is route.entry.go

4. Look at the generated routing file

```go
package main

import (
	"gin-annotation/_example/simple/controller"
	"gin-annotation/_example/simple/middleware"
	"github.com/gin-gonic/gin"
)

func Route(e *gin.Engine) {
	api := e.Group("/api", middleware.Log)
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/hello-world", middleware.Auth, controller.HelloWorld)
		}
	}
}
``` 

5. The last step , call Route() at main.go

```go
package main

import (
	"github.com/gin-gonic/gin"
	"path"
)

func main() {
	e := gin.Default()
	Route(e)
	_ = e.Run("ip:port")
}
```

## Annotations

- handlers
    - [groups](#groups-annotation)
    - [path](#path-annotation)
    - [method](#method-annotation)
    - [need](#need-annotation)
- middlewares
    - [id](#id-annotation)
    - [group](#group-annotation)
  
- [notice](#notice)

### Groups Annotation

Each groups-annotation consists of multiple groups separated by spaces.

### Path Annotation

The last element of path.

### Method Annotation

GET,POST,DELETE,PATCH,OPTIONS or ANY.

### Need Annotation

> need:id1 id2,

Each element of need-annotation is the id of the middleware.

### Id Annotation

The unique id of middleware.

### Group Annotation

```
/* example
[
  id:example,
  group:/api/v1/@1, 
  group:/api/v2/@1
]
*/
```

Each middleware can have multiple group-annotations,

The number behind @ is the priority of middleware.

### Notice
* Don't write extra "," in the last item of an annotation.
```
/* example
[
  id:example,
  group:/api/v1/@1    <- here
]
*/
```