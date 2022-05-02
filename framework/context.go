package framework

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	ctx            context.Context
	handlers       []ControllerHandler
	index          int // 当前请求调用到调用链的那个节点
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		ctx:            r.Context(),
		index:          -1,
	}
}

// base

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

func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context) QueryInt(key string, def int) int {
	hash := ctx.QueryAll()
	vals, ok := hash[key]
	if ok {
		if len(vals) > 0 {
			intval, err := strconv.Atoi(vals[len(vals)-1])
			if err != nil {
				return def
			}
			return intval
		}
	}

	return def
}

func (ctx *Context) QueryString(key string, def string) string {
	hash := ctx.QueryAll()
	vals, ok := hash[key]
	if ok {
		if len(vals) > 0 {
			return vals[len(vals)-1]
		}
	}
	return def
}

func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request != nil {
		// 强制类型转换
		return map[string][]string(ctx.request.URL.Query())
	}

	return map[string][]string{}
}

// response

func (ctx *Context) Json(status int, data interface{}) error {
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
