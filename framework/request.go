package framework

type IRequest interface {
	// 请求地址url中带参数
	// foo.com?a=1&b=2

	QueryInt(key string, def int) (int, bool)
	QueryString(key string, def string) (string, bool)

	// 路由匹配中带参数

	//ParamString(key string, def string) (string, bool)

	// form表单中携带参数

	//FormString(key string, def string) (string, bool)
}
