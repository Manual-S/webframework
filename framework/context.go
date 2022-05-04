package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cast"
)

type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	writeMux       *sync.Mutex
	ctx            context.Context
	hasTimeout     bool
	handlers       []ControllerHandler
	index          int               // 当前请求调用到调用链的那个节点
	params         map[string]string // 路由匹配的函数
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		ctx:            r.Context(),
		index:          -1,
		writeMux:       &sync.Mutex{},
	}
}

// base

func (ctx *Context) WriteMux() *sync.Mutex {
	return ctx.writeMux
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.responseWriter
}

func (ctx *Context) SetHasTimeout() {
	ctx.hasTimeout = true
}

func (ctx *Context) HasTimeout() bool {
	return ctx.hasTimeout
}

// context

func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	// todo
	return time.Time{}, false
}

func (ctx *Context) Err() error {
	// todo
	return nil
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

func (ctx *Context) Value(key interface{}) interface{} {
	// todo
	return nil
}

// request

func (ctx *Context) QueryInt(key string, def int) (int, bool) {
	hash := ctx.QueryAll()
	vals, ok := hash[key]
	if ok {
		if len(vals) > 0 {
			return cast.ToInt(vals[len(vals)-1]), true
		}
	}

	return def, false
}
func (ctx *Context) QueryInt64(key string, def int64) (int64, bool) {
	hash := ctx.QueryAll()
	vals, ok := hash[key]
	if ok {
		if len(vals) > 0 {
			return cast.ToInt64(vals[len(vals)-1]), true
		}
	}

	return def, false
}
func (ctx *Context) QueryFloat64(key string, def float64) (float64, bool) {
	hash := ctx.QueryAll()
	vals, ok := hash[key]
	if ok {
		if len(vals) > 0 {
			return cast.ToFloat64(vals[len(vals)-1]), true
		}
	}

	return def, false
}
func (ctx *Context) QueryFloat32(key string, def float32) (float32, bool) {
	hash := ctx.QueryAll()
	vals, ok := hash[key]
	if ok {
		if len(vals) > 0 {
			return cast.ToFloat32(vals[len(vals)-1]), true
		}
	}

	return def, false
}
func (ctx *Context) QueryBool(key string, def bool) (bool, bool) {
	hash := ctx.QueryAll()
	vals, ok := hash[key]
	if ok {
		if len(vals) > 0 {
			return cast.ToBool(vals[len(vals)-1]), true
		}
	}

	return def, false
}
func (ctx *Context) QueryStringSlice(key string, def []string) ([]string, bool) {
	hash := ctx.QueryAll()
	if vals, ok := hash[key]; ok {
		return vals, true
	}
	return def, false
}
func (ctx *Context) QueryString(key string, def string) (string, bool) {
	hash := ctx.QueryAll()
	vals, ok := hash[key]
	if ok {
		if len(vals) > 0 {
			return vals[len(vals)-1], true
		}
	}
	return def, false
}
func (ctx *Context) QueryArray(key string, def string) []string {
	// todo
	return nil
}
func (ctx *Context) Query(key string) interface{} {
	hash := ctx.QueryAll()
	if vals, ok := hash[key]; ok {
		return vals[0]
	}

	return nil
}
func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request != nil {
		// 强制类型转换
		return map[string][]string(ctx.request.URL.Query())
	}

	return map[string][]string{}
}

func (ctx *Context) ParamInt(key string, def int) (int, bool) {
	if val := ctx.Param(key); val != nil {
		return cast.ToInt(val), true
	}

	return def, false
}
func (ctx *Context) ParamInt64(key string, def int64) (int64, bool) {
	// todo
	return 0, false
}
func (ctx *Context) ParamFloat64(key string, def float64) (float64, bool) {
	// todo
	return 0, false
}
func (ctx *Context) ParamFloat32(key string, def float32) (float32, bool) {
	// todo
	return 0, false
}
func (ctx *Context) ParamBool(key string, def bool) (bool, bool) {
	// todo
	return false, false
}
func (ctx *Context) ParamString(key string, def string) (string, bool) {
	// todo
	return "", false
}
func (ctx *Context) Param(key string) interface{} {
	if ctx.params != nil {
		if val, ok := ctx.params[key]; ok {
			return val
		}
	}

	return nil
}

func (ctx *Context) BindJson(obj interface{}) error {
	if ctx.request != nil {
		body, err := ioutil.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}

		ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("request is empty")
	}
	return nil
}

func (ctx *Context) Uri() string {
	return ctx.request.RequestURI
}

func (ctx *Context) Method() string {
	return ctx.request.Method
}

func (ctx *Context) Host() string {
	return ctx.request.URL.Host
}

// ClientIp 获取ip地址
// todo 获取ip这里有需要http知识
// 参考资料 https://www.cnblogs.com/GaiHeiluKamei/p/13731791.html
func (ctx *Context) ClientIp() string {
	ip := ctx.request.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip
	}
	ip = ctx.request.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i
		}
	}
	if ip == "" {
		ip = ctx.request.RemoteAddr
	}
	return ip
}

// response

func (ctx *Context) Json(obj interface{}) IResponse {
	byt, err := json.Marshal(obj)
	if err != nil {
		return nil
	}
	ctx.SetHeader("Content-Type", "application/json")
	ctx.responseWriter.Write(byt)
	return nil
}

func (ctx *Context) SetHeader(key string, value string) IResponse {
	ctx.responseWriter.Header().Set(key, value)
	return nil
}

func (ctx *Context) SetCookie(key string, value string, maxAge int, path, domain string, secure, httpOnly bool) IResponse {
	// todo
	return nil
}

// SetStatus 设置状态码
func (ctx *Context) SetStatus(code int) IResponse {
	ctx.responseWriter.WriteHeader(code)
	return ctx
}

// SetOkStatus 设置200状态
func (ctx *Context) SetOkStatus() IResponse {
	ctx.responseWriter.WriteHeader(http.StatusOK)
	return ctx
}

//

func (ctx *Context) Next() error {
	ctx.index++
	if ctx.index < len(ctx.handlers) {
		// 有没有执行完的handler
		err := ctx.handlers[ctx.index](ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ctx *Context) SetHandlers(handlers []ControllerHandler) {
	ctx.handlers = handlers
}

func (ctx *Context) SetParams(params map[string]string) {
	ctx.params = params
}
