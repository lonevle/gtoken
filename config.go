package gtoken

import (
	"time"
)

type TokenStyle string

const (
	TokenStyleUUID    TokenStyle = "uuid"
	TokenStyleRand64  TokenStyle = "rand64"
	TokenStyleRand128 TokenStyle = "rand128"
)

type Config struct {
	CacheMode        int8          `c:"cacheMode"`        // 缓存模式 1 gcache 2 gredis 默认1
	CachePreKey      string        `c:"cachePreKey"`      // 缓存key前缀
	TokenStyle       TokenStyle    `c:"tokenStyle"`       // Token风格 默认 uuid
	TokenExpire      time.Duration `c:"tokenExpire"`      // Token过期时间 默认 2小时
	MultiLogin       bool          `c:"multiLogin"`       // 是否支持多端登录，默认false, （为true时支持多端登录，为false时仅允许一个端点登录）
	EnableRenew      bool          `c:"enableRenew"`      // 是否开启续期模式，默认false, （为true时在token过期前续期，为false时不续期）
	RenewThreshold   time.Duration `c:"renewThreshold"`   // 续期阈值，默认30分钟, 续签时间必须小于TokenExpire（距离过期时间多少开始续签）
	RenewTimeout     time.Duration `c:"renewTimeout"`     // 续期超时时间（默认3秒）
	AuthExcludePaths []string      `c:"authExcludePaths"` // 拦截排除地址 如: /login
	AuthHeaderName   string        `c:"authHeaderName"`   // 验证header的名称 默认 Authorization
	TokenPrefix      string        `c:"tokenPrefix"`      // token前缀 如`Bearer `
}

// 设置默认值
func (c *Config) setDefault() {
	if c.CacheMode == 0 {
		c.CacheMode = 1
	}
	if c.CachePreKey == "" {
		c.CachePreKey = "gtoken_"
	}
	if c.TokenStyle == "" {
		c.TokenStyle = TokenStyleUUID
	}
	if c.TokenExpire == 0 {
		c.TokenExpire = 2 * time.Hour // 默认2小时
	}
	if c.RenewThreshold == 0 {
		c.RenewThreshold = 30 * time.Minute // 默认30分钟
	}
	if c.RenewThreshold >= c.TokenExpire {
		c.RenewThreshold = c.TokenExpire / 2 // 续签时间必须小于TokenExpire
	}
	if c.RenewTimeout == 0 {
		c.RenewTimeout = 3 * time.Second
	}
	if c.AuthHeaderName == "" {
		c.AuthHeaderName = "Authorization"
	}
}
