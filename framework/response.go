package framework

type IResponse interface {
	Json(obj interface{}) IResponse

	SetHeader(key string, value string) IResponse

	SetCookie(key string, value string, maxAge int, path, domain string, secure, httpOnly bool) IResponse

	// 状态码相关

	// SetStatus 设置状态码
	SetStatus(code int) IResponse
	// SetOkStatus 设置200状态
	SetOkStatus() IResponse
}
