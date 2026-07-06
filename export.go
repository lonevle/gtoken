package gtoken

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Login 用户登录，生成并返回 token
// ctx: 上下文
// userID: 用户唯一标识
// return: token 和错误信息
func Login(ctx context.Context, userID string) (string, error) {
	if Instance == nil {
		return "", gerror.NewCode(gcode.CodeInternalError, MsgErrInstanceNotInit)
	}
	return Instance.Login(ctx, userID)
}

// Logout 用户退出登录，清除用户相关的 token 缓存
// ctx: 上下文
// userID: 用户唯一标识
// return: 错误信息
func Logout(ctx context.Context, userID string) error {
	if Instance == nil {
		return gerror.NewCode(gcode.CodeInternalError, MsgErrInstanceNotInit)
	}
	return Instance.Logout(ctx, userID)
}

// LogoutByToken 通过 token 退出登录，清除该 token 及关联的用户缓存
// ctx: 上下文
// token: 用户 token
// return: 错误信息
func LogoutByToken(ctx context.Context, token string) error {
	if Instance == nil {
		return gerror.NewCode(gcode.CodeInternalError, MsgErrInstanceNotInit)
	}
	return Instance.LogoutByToken(ctx, token)
}

// GetUserIDByToken 通过 token 获取用户 ID
// ctx: 上下文
// token: 用户 token
// return: 用户 ID 和错误信息
func GetUserIDByToken(ctx context.Context, token string) (*gvar.Var, error) {
	if Instance == nil {
		return nil, gerror.NewCode(gcode.CodeInternalError, MsgErrInstanceNotInit)
	}
	return Instance.GetUserIDByToken(ctx, token)
}

// GetTokenByUserKey 通过用户 ID 获取 token
// ctx: 上下文
// userID: 用户唯一标识
// return: token 和错误信息
func GetTokenByUserKey(ctx context.Context, userID string) (string, error) {
	if Instance == nil {
		return "", gerror.NewCode(gcode.CodeInternalError, MsgErrInstanceNotInit)
	}
	return Instance.GetTokenByUserKey(ctx, userID)
}

// IsLogin 验证 token 是否有效，判断用户是否已登录
// ctx: 上下文
// token: 用户 token
// return: 用户 ID、是否登录、错误信息
func IsLogin(ctx context.Context, token string) (userID *gvar.Var, isLoggedIn bool, err error) {
	if Instance == nil {
		return nil, false, gerror.NewCode(gcode.CodeInternalError, MsgErrInstanceNotInit)
	}
	return Instance.IsLogin(ctx, token)
}

// GenerateToken 生成新的 token
// userID: 用户唯一标识
// return: token 和错误信息
func GenerateToken(userID string) (string, error) {
	if Instance == nil {
		return "", gerror.NewCode(gcode.CodeInternalError, MsgErrInstanceNotInit)
	}
	return Instance.GenerateToken(userID)
}

// AddExcludePaths 添加认证排除路径，这些路径不需要 token 验证
// urls: 需要排除的路径列表
func AddExcludePaths(urls ...string) {
	if Instance != nil {
		Instance.AddExcludePaths(urls...)
	}
}

// SearchPath 检查路径是否在排除列表中
// url: 要检查的路径
// return: 是否在排除列表中
func SearchPath(url string) bool {
	if Instance == nil {
		return false
	}
	return Instance.SearchPath(url)
}

// RemovePaths 移除排除路径
// urls: 要移除的路径列表
func RemovePaths(urls ...string) {
	if Instance != nil {
		Instance.RemovePaths(urls...)
	}
}

// ClearPaths 清空所有排除路径
func ClearPaths() {
	if Instance != nil {
		Instance.ClearPaths()
	}
}
