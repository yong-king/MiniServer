package errorx


const(
	DefaultErrorCode = 1001
	RpcErroCode = 1002
	SqlErrorCode = 1003
	QuerNoFoundErrorCode = 1004
	RedisErrorCode = 1005
)

// CodeError 自定义错误类型
type CodeError struct{
	Code int `json:"code"`
	Msg string `json:"msg"`
}

// CodeErrorResponse 自定义响应错误类型
type CodeErrorResponse struct{
	Code int `josn:"code"`
	Msg string `json:"msg"`
}

// NewCodeError 返回自定义错误
func NewCodeError(code int, msg string) error {
	return CodeError{
		Code: code,
		Msg: msg,
	}
}
// Error CodeError实现error接口
func (e CodeError) Error() string {
	return e.Msg
}

// NewDefaultCodeError 返回默认自定义错误
func NewDefaultCodeError(msg string) error {
	return CodeError{
		Code: DefaultErrorCode,
		Msg: msg,
	}
}

// Data 返回自定义类型的错误响应
func (e *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code: e.Code,
		Msg: e.Msg,
	}
}