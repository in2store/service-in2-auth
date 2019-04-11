package errors

import (
	"github.com/johnnyeven/libtools/courier/status_error"
)

func init() {
	status_error.StatusErrorCodes.Register("BadRequest", 400100000, "请求参数错误", "", false)
	status_error.StatusErrorCodes.Register("Unauthorized", 401100000, "未授权", "", true)
	status_error.StatusErrorCodes.Register("Forbidden", 403100000, "不允许操作", "", true)
	status_error.StatusErrorCodes.Register("NotFound", 404100000, "未找到", "", false)
	status_error.StatusErrorCodes.Register("ChannelNotFound", 404100001, "Git托管通道配置未找到", "", true)
	status_error.StatusErrorCodes.Register("Conflict", 409100000, "操作冲突", "", true)
	status_error.StatusErrorCodes.Register("InternalError", 500100000, "内部处理错误", "", false)
}
