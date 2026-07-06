package gtoken

import (
	"context"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gmlock"
)

type Manager struct {
	config         *Config
	prefix         string
	cache          *cache
	excludePaths   *gmap.Map
	tokenPrefixLen int // token前缀长度
}

// 登录
func (m *Manager) Login(ctx context.Context, userID string) (string, error) {
	gmlock.Lock(userID)
	defer gmlock.Unlock(userID)

	userKey := m.getUserKey(userID)

	// 先从缓存中获取token
	v, err := m.cache.Get(ctx, userKey)
	if err != nil {
		g.Log().Warning(ctx, "[GToken] get user cache error", userKey, err)
	} else if !v.IsNil() {
		// 如果存在token
		if m.config.MultiLogin {
			return v.String(), nil
		}
		// 把旧token删除
		tokenKey := m.getTokenKey(v.String())
		if err := m.cache.Remove(ctx, tokenKey, userKey); err != nil {
			g.Log().Warning(ctx, "[GToken] remove old token failed", userKey, err)
		}
	}

	// 不允许多端登录或者缓存中没有token，则生成新的token
	token, err := m.GenerateToken(userID)
	if err != nil {
		return "", err
	}

	// 设置token到缓存中
	tokenKey := m.getTokenKey(token)

	// 设置token对应的用户信息到缓存中
	err = m.cache.Set(ctx, tokenKey, userID, m.config.TokenExpire)
	if err != nil {
		return "", err
	}
	// 设置用户对应的token到缓存中
	err = m.cache.Set(ctx, userKey, token, m.config.TokenExpire)
	if err != nil {
		// 回滚已写入的tokenKey，避免脏数据
		m.cache.Remove(ctx, tokenKey)
		return "", err
	}

	return token, nil
}

// 退出登录
func (m *Manager) Logout(ctx context.Context, userID string) error {
	gmlock.Lock(userID)
	defer gmlock.Unlock(userID)

	token, err := m.GetTokenByUserKey(ctx, userID)
	if err != nil {
		return err
	}
	return m.cache.Remove(ctx, m.getTokenKey(token), m.getUserKey(userID))
}

// 通过token退出登录
func (m *Manager) LogoutByToken(ctx context.Context, token string) error {
	userID, err := m.GetUserIDByToken(ctx, token)
	if err != nil {
		return err
	}
	userIDStr := userID.String()
	gmlock.Lock(userIDStr)
	defer gmlock.Unlock(userIDStr)

	return m.cache.Remove(ctx, m.getTokenKey(token), m.getUserKey(userIDStr))
}

// 通过token获取userID, UserID过期或不存在error
func (m *Manager) GetUserIDByToken(ctx context.Context, token string) (*gvar.Var, error) {
	if token == "" {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, MsgErrTokenEmpty)
	}
	tokenKey := m.getTokenKey(token)
	v, err := m.cache.Get(ctx, tokenKey)
	if err != nil {
		return nil, err
	}
	if v.IsNil() {
		return nil, gerror.NewCode(gcode.CodeNotFound, MsgErrExpiredOrNotExist)
	}
	return v, nil
}

// 通过userID获取token， token过期或不存在error
func (m *Manager) GetTokenByUserKey(ctx context.Context, userID string) (string, error) {
	if userID == "" {
		return "", gerror.NewCode(gcode.CodeMissingParameter, MsgErrUserIDEmpty)
	}
	userKey := m.getUserKey(userID)
	v, err := m.cache.Get(ctx, userKey)
	if err != nil {
		return "", err
	}
	if v.IsNil() {
		return "", gerror.NewCode(gcode.CodeNotFound, MsgErrExpiredOrNotExist)
	}
	return v.String(), nil
}

func (m *Manager) IsLogin(ctx context.Context, token string) (userID *gvar.Var, isLoggedIn bool, err error) {
	if token == "" {
		return nil, false, gerror.NewCode(gcode.CodeMissingParameter, MsgErrTokenEmpty)
	}
	// 通过token能获取到user
	userID, err = m.GetUserIDByToken(ctx, token)
	if err != nil || userID == nil {
		return nil, false, err
	}
	userKey := m.getUserKey(userID.String())
	// 并且user存在
	ok, _ := m.cache.Exists(ctx, userKey)
	if !ok {
		return nil, false, gerror.NewCode(gcode.CodeNotFound, MsgErrExpiredOrNotExist)
	}

	// 如果开启续期模式，则在token过期前续期
	if m.config.EnableRenew {
		go m.renewCache(m.getTokenKey(token), userKey)
	}

	return userID, true, nil
}

// 缓存续期
func (m *Manager) renewCache(keys ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), m.config.RenewTimeout)
	defer cancel()
	for _, v := range keys {
		expire, err := m.cache.GetExpire(ctx, v)
		if err != nil || expire <= 0 || expire >= m.config.RenewThreshold {
			continue
		}
		success, err := m.cache.Renew(ctx, v, m.config.TokenExpire)
		if err != nil || !success {
			g.Log().Warning(ctx, "[GToken] renew cache failed", v)
		}
	}
}

// 生成token
func (m *Manager) GenerateToken(userID string) (string, error) {
	if userID == "" {
		return "", gerror.NewCode(gcode.CodeMissingParameter, MsgErrUserIDEmpty)
	}

	switch m.config.TokenStyle {
	case TokenStyleUUID:
		return generateUUIDToken(userID), nil
	case TokenStyleRand64:
		return generateRandToken(64), nil
	case TokenStyleRand128:
		return generateRandToken(128), nil
	default:
		return generateUUIDToken(userID), nil
	}
}

// 获取用户缓存key
func (m *Manager) getUserKey(userID string) string {
	return m.prefix + UserKeyPrefix + userID
}

// 获取token缓存key
func (m *Manager) getTokenKey(tokenValue string) string {
	return m.prefix + TokenKeyPrefix + tokenValue
}

// 添加排除路径
func (m *Manager) AddExcludePaths(urls ...string) {
	for _, url := range urls {
		m.excludePaths.Set(url, 1)
	}
}

// 检查是否包含该路径
func (m *Manager) SearchPath(url string) bool {
	_, found := m.excludePaths.Search(url)
	return found
}

// 移除排除路径
func (m *Manager) RemovePaths(url ...string) {
	arrUrl := make([]any, len(url))
	for k, v := range url {
		arrUrl[k] = v
	}
	m.excludePaths.Removes(arrUrl)
}

// 清空排除路径
func (m *Manager) ClearPaths() {
	m.excludePaths.Clear()
}
