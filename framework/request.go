package framework

type IRequest interface {
	// 请求地址url中带参数
	// foo.com?a=1&b=2

	QueryInt(key string, def int) (int, bool)
	QueryInt64(key string, def int64) (int64, bool)
	QueryFloat64(key string, def float64) (float64, bool)
	QueryFloat32(key string, def float32) (float32, bool)
	QueryBool(key string, def bool) (bool, bool)
	QueryString(key string, def string) (string, bool)
	QueryStringSlice(key string, def []string) ([]string, bool)
	Query(key string) interface{}

	// 路由匹配中带参数

	ParamInt(key string, def int) (int, bool)
	ParamInt64(key string, def int64) (int64, bool)
	ParamFloat64(key string, def float64) (float64, bool)
	ParamFloat32(key string, def float32) (float32, bool)
	ParamBool(key string, def bool) (bool, bool)
	ParamString(key string, def string) (string, bool)
	Param(key string) interface{}

	// form表单中携带参数

	//FormString(key string, def string) (string, bool)

	// bindXXX函数

	BindJson(obj interface{}) error

	// 基础信息的封装

	Uri() string
	Method() string
	Host() string
	ClientIp() string

	// header

}
