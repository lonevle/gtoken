package gtoken

import (
	"github.com/gogf/gf/v2/util/grand"
	"github.com/gogf/gf/v2/util/guid"
)

// 生成UUID token
func generateUUIDToken(userID string) string {
	return guid.S([]byte(userID))
}

// 生成随机token
func generateRandToken(length int) string {
	return grand.Letters(length)
}
