package errors

import (
	"net/http"

	"github.com/johnnyeven/libtools/courier/status_error"
)

//go:generate libtools gen error
const ServiceStatusErrorCode = 0 * 1e3 // todo rename this

const (
	// 请求参数错误
	BadRequest status_error.StatusErrorCode = http.StatusBadRequest*1e6 + ServiceStatusErrorCode + iota
)

const (
	// 未找到
	NotFound status_error.StatusErrorCode = http.StatusNotFound*1e6 + ServiceStatusErrorCode + iota
)

const (
	// @errTalk 未授权
	Unauthorized status_error.StatusErrorCode = http.StatusUnauthorized*1e6 + ServiceStatusErrorCode + iota
	// @errTalk 用户名或密码错误
	UserNameOrPasswordError
	// @errTalk 无效的身份令牌
	InvalidToken
)

const (
	// @errTalk 操作冲突
	Conflict status_error.StatusErrorCode = http.StatusConflict*1e6 + ServiceStatusErrorCode + iota
	// 用户已存在
	UserAlreadyExist
)

const (
	// @errTalk 不允许操作
	Forbidden status_error.StatusErrorCode = http.StatusForbidden*1e6 + ServiceStatusErrorCode + iota
	// 参数不完整
	ArgumentsNotEnough
)

const (
	// 内部处理错误
	InternalError status_error.StatusErrorCode = http.StatusInternalServerError*1e6 + ServiceStatusErrorCode + iota
	// 数据库操作内部错误
	InternalDbError
)
