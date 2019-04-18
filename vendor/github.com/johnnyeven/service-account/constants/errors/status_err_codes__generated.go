package errors

import (
	"github.com/johnnyeven/libtools/courier/status_error"
)

func init() {
	status_error.StatusErrorCodes.Register("BadRequest", 400000000, "请求参数错误", "", false)
	status_error.StatusErrorCodes.Register("Unauthorized", 401000000, "未授权", "", true)
	status_error.StatusErrorCodes.Register("UserNameOrPasswordError", 401000001, "用户名或密码错误", "", true)
	status_error.StatusErrorCodes.Register("InvalidToken", 401000002, "无效的身份令牌", "", true)
	status_error.StatusErrorCodes.Register("Forbidden", 403000000, "不允许操作", "", true)
	status_error.StatusErrorCodes.Register("ArgumentsNotEnough", 403000001, "参数不完整", "", false)
	status_error.StatusErrorCodes.Register("NotFound", 404000000, "未找到", "", false)
	status_error.StatusErrorCodes.Register("Conflict", 409000000, "操作冲突", "", true)
	status_error.StatusErrorCodes.Register("UserAlreadyExist", 409000001, "用户已存在", "", false)
	status_error.StatusErrorCodes.Register("InternalError", 500000000, "内部处理错误", "", false)
	status_error.StatusErrorCodes.Register("InternalDbError", 500000001, "数据库操作内部错误", "", false)
}
