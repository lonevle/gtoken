package gtoken

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// 全局实例
var (
	Instance *Manager
	once     sync.Once
)

// 初始化GToken实例
func Init(config *Config) *Manager {
	once.Do(func() {
		config.setDefault()

		Instance = &Manager{
			config:         config,
			prefix:         config.CachePreKey,
			cache:          newCache(config.CacheMode),
			excludePaths:   gmap.New(true),
			tokenPrefixLen: len(config.TokenPrefix),
		}
		Instance.AddExcludePaths(config.AuthExcludePaths...)
		g.Log().Debug(gctx.New(), "[GToken]initialized", config)
	})
	return Instance
}

// 从配置文件中读取
func NewConfigFromCtx(ctx context.Context) *Config {
	var config *Config
	err := g.Cfg().MustGet(ctx, "gToken").Struct(&config)
	if err != nil {
		panic("gToken config init fail: " + err.Error())
	}
	if config == nil {
		panic("gToken config not configured")
	}
	return config
}
