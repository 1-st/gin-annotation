# gin-annotation

[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/1-st/gin-annotation/main/LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/1-st/gin-annotation)
[![Go](https://github.com/1-st/gin-annotation/workflows/Go/badge.svg?branch=main)](https://github.com/1-st/gin-annotation/actions)

一个实现gin框架注解路由的命令行工具

 <img src="https://raw.githubusercontent.com/1-st/logos/master/gin-annotation/logo.png" width = "50%" alt="logo" align=center /> 


## 特性

* 通过操作golang [AST](https://en.wikipedia.org/wiki/Abstract_syntax_tree) 进行代码生成
* 组路由支持
* 通过注解实现中间件排序
* 非侵入式设计

## 快速开始

1. 安装.

```shell
go get github.com/1-st/gin-annotation
```

2. 在任何文件编写你的handler函数.

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

3. 在项目目录执行gin-annotation ./: (例如: *_example/simple* ; 你可以指定多个目录)

```sh
$ gin-annotation ./
$ ls
controller main.go route.entry.go(!!!new file)
```

> 提示: 新文件的名字由环境变量 <font color=#ee00ee>GIN_ANNOTATION_FILE</font> 决定
, 默认是route.entry.go

4. 查看生成的route.entry.go

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

5. 最后一步，在main函数调用Route()

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

## 注解

- handler
    - [groups](#groups-annotation)
    - [path](#path-annotation)
    - [method](#method-annotation)
    - [need](#need-annotation)
- middlewares
    - [id](#id-annotation)
    - [group](#group-annotation)

- [notice](#notice)

### Groups Annotation

每个路由组之间由空格分隔开

### Path Annotation

路径的最后一个元素

### Method Annotation

GET,POST,DELETE,PATCH,OPTIONS or ANY.

### Need Annotation

> need:id1 id2,

每个元素都是middleware的ID

### Id Annotation

middleware的唯一ID.

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

每个middleware可以拥有多个group-annotation,

@之后的数字是middleware在group中的优先级.

### Notice
* 在注解的最后注意不要加上多余的','
```
/* example
[
  id:example,
  group:/api/v1/@1    <- here
]
*/
```
