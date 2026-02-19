# EX

EX是一个基于Go标准库`net/http`扩展的轻量级Web框架。

## 特性

- 🚀 轻量级，仅依赖Go标准库
- 📦 支持路由分组，静态资源服务，Websocket和SSE
- 🔧 中间件支持
- 🛡️ 内置Logger，Recovery，跨域，Request ID，JWT等中间件
- 🎯 简洁的API设计

## 安装

```bash
go get github.com/seaung/ex
```

## 快速开始

### 基础示例

```go
package main

import (
    "github.com/seaung/ex"
)

func main() {
    engine := ex.NewEngine()
    
    engine.GET("/hello", func(ctx *ex.Context) {
        ctx.String(200, "Hello World!")
    })
    
    engine.Run(":9527")
}
```

### 使用默认引擎（包含内置中间件）

```go
package main

import (
    "github.com/seaung/ex"
)

func main() {
    engine := ex.DefaultEngine()
    
    engine.GET("/hello", func(ctx *ex.Context) {
        ctx.String(200, "Hello World!")
    })
    
    engine.Run(":9527")
}
```

## 路由

### HTTP 方法

框架支持常见的 HTTP 方法：

```go
engine.GET("/users", listUsers)
engine.POST("/users", createUser)
engine.PUT("/users/:id", updateUser)
engine.DELETE("/users/:id", deleteUser)
```

### 路由分组

使用路由分组可以更好地组织 API 结构，并为一组路由统一添加中间件：

```go
func main() {
    engine := ex.NewEngine()
    
    api := engine.AddGroup("/api")
    api.GET("/hello", func(ctx *ex.Context) {
        ctx.String(200, "API Hello")
    })
    
    v1 := api.AddGroup("/v1")
    v1.GET("/users", func(ctx *ex.Context) {
        ctx.String(200, "v1 users")
    })
    
    engine.Run(":9527")
}
```

## 中间件

### 注册中间件

```go
func main() {
    engine := ex.NewEngine()
    
    engine.Use(Logger(), Recovery())
    
    engine.GET("/hello", func(ctx *ex.Context) {
        ctx.String(200, "Hello World!")
    })
    
    engine.Run(":9527")
}
```

### 自定义中间件

```go
func Auth() ex.HandlerFunc {
    return func(ctx *ex.Context) {
        token := ctx.Query("token")
        if token == "" {
            ctx.String(401, "Unauthorized")
            ctx.Abort()
            return
        }
        ctx.Next()
    }
}

func main() {
    engine := ex.NewEngine()
    
    api := engine.AddGroup("/api")
    api.Use(Auth())
    api.GET("/profile", func(ctx *ex.Context) {
        ctx.String(200, "Profile Data")
    })
    
    engine.Run(":9527")
}
```

### 中间件执行顺序

中间件按照注册顺序执行：

```go
engine.Use(
    func(ctx *ex.Context) {
        fmt.Println("Middleware 1 - Before")
        ctx.Next()
        fmt.Println("Middleware 1 - After")
    },
    func(ctx *ex.Context) {
        fmt.Println("Middleware 2 - Before")
        ctx.Next()
        fmt.Println("Middleware 2 - After")
    },
)
```

输出顺序：
```
Middleware 1 - Before
Middleware 2 - Before
Middleware 2 - After
Middleware 1 - After
```

## Context

`Context` 封装了请求和响应的上下文信息：

```go
type Context struct {
    Writer     http.ResponseWriter
    Req        *http.Request
    Path       string
    Method     string
    StatusCode int
}
```

### Context 方法

| 方法 | 说明 |
|------|------|
| `Query(key string) string` | 获取 URL 查询参数 |
| `String(code int, msg string)` | 返回字符串响应 |
| `Status(code int)` | 设置响应状态码 |
| `Next()` | 执行下一个中间件 |
| `Abort()` | 终止中间件链 |

### 示例

```go
engine.GET("/search", func(ctx *ex.Context) {
    keyword := ctx.Query("q")
    ctx.String(200, "Search: "+keyword)
})
```

## 内置中间件

### Logger

记录请求日志：

```go
engine.Use(ex.Logger())
```

### Recovery

恢复 panic，防止服务崩溃：

```go
engine.Use(ex.Recovery())
```

## API 参考

### Engine

| 方法 | 说明 |
|------|------|
| `NewEngine() *Engine` | 创建一个新的引擎实例 |
| `DefaultEngine() *Engine` | 创建一个带有 Logger 和 Recovery 中间件的引擎 |
| `Run(addr string) error` | 启动 HTTP 服务器 |
| `GET(path string, handlers ...HandlerFunc)` | 注册 GET 路由 |
| `POST(path string, handlers ...HandlerFunc)` | 注册 POST 路由 |
| `PUT(path string, handlers ...HandlerFunc)` | 注册 PUT 路由 |
| `DELETE(path string, handlers ...HandlerFunc)` | 注册 DELETE 路由 |
| `Use(middlewares ...HandlerFunc)` | 注册全局中间件 |
| `AddGroup(prefix string) *RouterGroup` | 创建路由分组 |

### RouterGroup

| 方法 | 说明 |
|------|------|
| `GET(path string, handlers ...HandlerFunc)` | 注册 GET 路由 |
| `POST(path string, handlers ...HandlerFunc)` | 注册 POST 路由 |
| `PUT(path string, handlers ...HandlerFunc)` | 注册 PUT 路由 |
| `DELETE(path string, handlers ...HandlerFunc)` | 注册 DELETE 路由 |
| `Use(middlewares ...HandlerFunc)` | 注册分组级别中间件 |
| `AddGroup(prefix string) *RouterGroup` | 创建子分组 |

## 示例项目

查看 `examples/` 目录获取更多示例：

- [basic](./examples/basic) - 基础用法
- [middlewares](./examples/middlewares) - 中间件使用
- [controller](./examples/controller) - Controller使用

---
that's all
