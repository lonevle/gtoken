package gtoken

// 缓存模式
const (
	CacheModeCache = 1
	CacheModeRedis = 2
)

// 常量
const (
	UserKeyPrefix  = "account:"
	TokenKeyPrefix = "token:"
)

// 错误信息
const (
	MsgErrUserIDEmpty       = "userID empty"
	MsgErrDataEmpty         = "cache value is nil"
	MsgErrTokenEmpty        = "token empty"
	MsgErrExpiredOrNotExist = "Data Expired or not Exist"
)
