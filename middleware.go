package gtoken

import (
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

type Middleware struct {
	ResFunc func(r *ghttp.Request, err error)
}

func NewMiddleware(resFunc ...func(r *ghttp.Request, err error)) *Middleware {
	m := &Middleware{}
	if len(resFunc) > 0 {
		m.ResFunc = resFunc[0]
	} else {
		m.ResFunc = func(r *ghttp.Request, err error) {
			r.Response.WriteJson(ghttp.DefaultHandlerResponse{
				Code:    gcode.CodeBusinessValidationFailed.Code(),
				Message: err.Error(),
			})
		}
	}
	return m
}

func (m *Middleware) Auth(r *ghttp.Request) {
	if m.hasAccessExcludePath(r) {
		r.Middleware.Next()
		return
	}
	token, err := GetRequestToken(r)
	if err != nil {
		m.ResFunc(r, err)
		return
	}
	// 检查token是否有效
	userID, ok, err := Instance.IsLogin(r.Context(), token)
	if err != nil || !ok {
		m.ResFunc(r, err)
		return
	}
	r.SetCtxVar("userID", userID)
	r.Middleware.Next()
}

// 是否在排除路径
func (m *Middleware) hasAccessExcludePath(r *ghttp.Request) bool {
	handler := r.GetServeHandler()
	if handler == nil || handler.Handler.Router == nil {
		return false
	}
	route := handler.Handler.Router.Uri
	return Instance.SearchPath(route)
}

func GetRequestToken(r *ghttp.Request) (string, error) {
	authHeader := r.Header.Get(Instance.config.AuthHeaderName)
	if authHeader == "" {
		return "", gerror.NewCode(gcode.CodeMissingParameter, MsgErrTokenEmpty)
	}

	// 只保留token，去除前缀
	if Instance.tokenPrefixLen > 0 {
		authHeader = strings.TrimPrefix(authHeader, Instance.config.TokenPrefix)
	}
	return authHeader, nil
}
