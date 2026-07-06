package gtoken

import (
	"context"
	"errors"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
)

type cache struct {
	cache *gcache.Cache
}

func newCache(mode int8) *cache {
	c := &cache{
		cache: gcache.New(),
	}

	if mode == CacheModeRedis {
		c.cache.SetAdapter(gcache.NewAdapterRedis(g.Redis()))
	}

	return c
}

// 设置缓存
func (c *cache) Set(ctx context.Context, cacheKey string, cacheValue any, expire time.Duration) error {
	if cacheValue == nil {
		return errors.New(MsgErrDataEmpty)
	}

	return c.cache.Set(ctx, cacheKey, cacheValue, expire)
}

// 获取缓存
func (c *cache) Get(ctx context.Context, cacheKey string) (*gvar.Var, error) {
	dataVar, err := c.cache.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	return dataVar, nil
}

// 删除缓存, 传入与设置数据类型要一致，否则可能导致问题
func (c *cache) Remove(ctx context.Context, cacheKeys ...any) error {
	_, err := c.cache.Remove(ctx, cacheKeys...)
	return err
}

// 判断缓存是否存在
func (c *cache) Exists(ctx context.Context, cacheKey string) (bool, error) {
	exists, err := c.cache.Contains(ctx, cacheKey)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// 续期
func (c *cache) Renew(ctx context.Context, cacheKey string, expire time.Duration) (bool, error) {
	old, err := c.cache.UpdateExpire(ctx, cacheKey, expire)
	if err != nil {
		return false, err
	}
	// -1 表示续期失败
	if old == -1 {
		return false, nil
	}
	return true, nil
}

// 获取缓存剩余时间
func (c *cache) GetExpire(ctx context.Context, cacheKey string) (time.Duration, error) {
	expire, err := c.cache.GetExpire(ctx, cacheKey)
	if err != nil {
		return 0, err
	}
	return expire, nil
}
