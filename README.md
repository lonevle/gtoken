# gtoken

基于 GoFrame v2 框架的轻量级 Token 认证插件，提供登录、登出、token 续期等功能。

## 特性

- 支持多种缓存模式（gcache / gredis）
- 支持多种 Token 生成方式（UUID / 随机字符串）
- 单点登录与多端登录模式
- Token 自动续期机制
- HTTP Middleware 集成
- 排除路径配置

## 安装

```bash
go get github.com/lonevle/gtoken
```

## 快速开始

### 方式一：代码配置

```go
package main

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
    "github.com/lonevle/gtoken"
)

func main() {
    s := g.Server()
    
    // 初始化 gtoken
    gtoken.Init(&gtoken.Config{
        CacheMode:        gtoken.CacheModeCache, // 缓存模式：1 gcache 2 gredis
        CachePreKey:      "gtoken_",             // 缓存 key 前缀
        TokenStyle:       gtoken.TokenStyleUUID, // Token 风格
        TokenExpire:      2 * time.Hour,         // Token 过期时间
        MultiLogin:       false,                 // 是否支持多端登录
        EnableRenew:      true,                  // 是否开启续期
        RenewThreshold:   30 * time.Minute,      // 续期阈值
        AuthExcludePaths: []string{"/login"},    // 排除路径
        AuthHeaderName:   "Authorization",       // 认证 header 名称
        TokenPrefix:      "Bearer ",             // Token 前缀
    })
    
    // 注册中间件
    s.Group("/", func(group *ghttp.RouterGroup) {
        group.Middleware(gtoken.NewMiddleware().Auth)
        group.GET("/login", loginHandler)
        group.GET("/logout", logoutHandler)
        group.GET("/user", userHandler)
    })
    
    s.Run()
}

func loginHandler(r *ghttp.Request) {
    token, err := gtoken.Login(r.Context(), "user001")
    if err != nil {
        r.Response.WriteJsonExit(map[string]any{"code": 1, "msg": err.Error()})
    }
    r.Response.WriteJsonExit(map[string]any{"code": 0, "token": token})
}

func logoutHandler(r *ghttp.Request) {
    err := gtoken.Logout(r.Context(), "user001")
    if err != nil {
        r.Response.WriteJsonExit(map[string]any{"code": 1, "msg": err.Error()})
    }
    r.Response.WriteJsonExit(map[string]any{"code": 0, "msg": "success"})
}

func userHandler(r *ghttp.Request) {
    userID := r.GetCtxVar("userID").String()
    r.Response.WriteJsonExit(map[string]any{"code": 0, "userID": userID})
}
```

### 方式二：配置文件

在 `config.yaml` 中配置：

```yaml
gToken:
  cacheMode: 1              # 缓存模式：1 gcache 2 gredis
  cachePreKey: "gtoken_"    # 缓存 key 前缀
  tokenStyle: "uuid"        # Token 风格：uuid / rand64 / rand128
  tokenExpire: "2h"         # Token 过期时间，默认 2h（支持格式：2h, 30m, 7200s）
  multiLogin: false         # 是否支持多端登录
  enableRenew: true         # 是否开启续期，默认 false
  renewThreshold: "30m"     # 续期阈值，默认 30m（支持格式：30m, 1800s）
  renewTimeout: "3s"        # 续期超时时间，默认 3s
  authExcludePaths:         # 排除路径
    - /login
    - /register
  authHeaderName: "Authorization"  # 认证 header 名称
  tokenPrefix: "Bearer "           # Token 前缀
```

然后在代码中加载：

```go
func main() {
    config := gtoken.NewConfigFromCtx(gctx.New())
    gtoken.Init(config)
    // ...
}
```

## 配置说明

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `CacheMode` | `int8` | 1 | 缓存模式：1 gcache（内存），2 gredis（Redis） |
| `CachePreKey` | `string` | `gtoken_` | 缓存 key 前缀 |
| `TokenStyle` | `string` | `uuid` | Token 风格：`uuid`、`rand64`、`rand128` |
| `TokenExpire` | `time.Duration` | 2小时 | Token 过期时间 |
| `MultiLogin` | `bool` | `false` | 是否支持多端登录，true 允许多端，false 单点登录 |
| `EnableRenew` | `bool` | `false` | 是否开启续期模式 |
| `RenewThreshold` | `time.Duration` | 30分钟 | 续期阈值，需小于 TokenExpire |
| `RenewTimeout` | `time.Duration` | 3秒 | 续期操作超时时间 |
| `AuthExcludePaths` | `[]string` | `nil` | 认证排除路径，如 `/login` |
| `AuthHeaderName` | `string` | `Authorization` | 认证 header 名称 |
| `TokenPrefix` | `string` | `` | Token 前缀，如 `Bearer ` |

## API 方法

包级函数使用方式（推荐）：

### 登录

```go
token, err := gtoken.Login(ctx, userID string) (string, error)
```

### 登出

```go
err := gtoken.Logout(ctx, userID string) error
```

### 通过 Token 登出

```go
err := gtoken.LogoutByToken(ctx, token string) error
```

### 验证登录状态

```go
userID, isLoggedIn, err := gtoken.IsLogin(ctx, token string) (*gvar.Var, bool, error)
```

### 通过 Token 获取用户 ID

```go
userID, err := gtoken.GetUserIDByToken(ctx, token string) (*gvar.Var, error)
```

### 通过用户 ID 获取 Token

```go
token, err := gtoken.GetTokenByUserKey(ctx, userID string) (string, error)
```

### 生成 Token

```go
token, err := gtoken.GenerateToken(userID string) (string, error)
```

### 路径管理

```go
// 添加排除路径
gtoken.AddExcludePaths("/login", "/register")

// 检查路径是否在排除列表
exists := gtoken.SearchPath("/login")

// 移除排除路径
gtoken.RemovePaths("/login")

// 清空排除路径
gtoken.ClearPaths()
```

## Middleware 使用

### 默认响应

```go
group.Middleware(gtoken.NewMiddleware().Auth)
```

默认响应格式：

```json
{
    "code": 40000,
    "message": "error message",
    "data": null
}
```

### 自定义响应

```go
group.Middleware(gtoken.NewMiddleware(func(r *ghttp.Request, err error) {
    r.Response.WriteJsonExit(map[string]any{
        "code":    401,
        "message": "未登录或登录已过期",
    })
}).Auth)
```

## Token 风格

| 风格 | 说明 |
|------|------|
| `TokenStyleUUID` | 使用 UUID 生成，32 字节，包含 MAC、PID、时间戳、序列号、随机数 |
| `TokenStyleRand64` | 使用随机字符串生成，64 字节 |
| `TokenStyleRand128` | 使用随机字符串生成，128 字节 |

## 缓存模式

| 模式 | 说明 |
|------|------|
| `CacheModeCache` | 使用 GoFrame 内置缓存（gcache），适用于单实例部署 |
| `CacheModeRedis` | 使用 Redis 缓存，适用于多实例部署（需要提前配置 Redis） |

## 注意事项

1. **分布式锁**：当前使用 `gmlock`（进程内锁），在 Redis 多实例环境下无法提供分布式锁保护
2. **Token 续期**：`EnableRenew` 默认关闭（`false`），需显式开启
3. **单点登录**：`MultiLogin` 为 `false` 时，每次登录会生成新 Token 并失效旧 Token