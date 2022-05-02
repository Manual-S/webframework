package framework

import (
	"context"
	"encoding/json"
	"net/http"
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
	index          int // 当前请求调用到调用链的那个节点
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

func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request != nil {
		// 强制类型转换
		return map[string][]string(ctx.request.URL.Query())
	}

	return map[string][]string{}
}

func (ctx *Context) Formint(key string, def int) int {
	// todo
	return 0
}

func (ctx *Context) FormString(key string, def string) string {
	// todo
	return ""
}

func (ctx *Context) FormArray(key string, def []string) []string {
	// todo
	return nil
}

func (ctx *Context) FormAll() map[string][]string {
	// todo
	return nil
}
func (ctx *Context) BindJson(obj interface{}) error {
	// todo
	return nil
}

// response

func (ctx *Context) Json(status int, data interface{}) error {
	if ctx.HasTimeout() {
		return nil
	}
	ctx.responseWriter.Header().Set("Content-Type", "application/json")
	ctx.responseWriter.WriteHeader(status)
	byt, err := json.Marshal(data)
	if err != nil {
		ctx.responseWriter.WriteHeader(http.StatusInternalServerError)
		return err
	}
	ctx.responseWriter.Write(byt)
	return nil
}

func (ctx *Context) HTML() error {
	return nil
}

func (ctx *Context) Text() error {
	return nil
}

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
