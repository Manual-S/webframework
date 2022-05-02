package framework

import (
	"log"
	"net/http"
	"strings"
)

type Core struct {
	router       map[string]*Tree
	middleswares []ControllerHandler
}

func NewCore() *Core {
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()

	return &Core{
		router: router,
	}
}

func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request, response)

	router := c.FindRouteByRequest(request)
	if router == nil {
		ctx.Json(http.StatusNotFound, "not found")
		return
	}

	// 设置context中的handlers字段
	ctx.SetHandlers(router)

	err := ctx.Next()
	if err != nil {
		ctx.Json(http.StatusInternalServerError, "inner error")
		return
	}
}

func (c *Core) Get(url string, handler ...ControllerHandler) {
	// 将core中已经注册的中间件和get自己要注册的中间件结合起来
	allHandlers := append(c.middleswares, handler...)
	err := c.router["GET"].AddRouter(url, allHandlers...)
	if err != nil {
		log.Fatal("add router error" + url)
	}
}

func (c *Core) Post(url string, handler ControllerHandler) {

}

func (c *Core) FindRouteByRequest(req *http.Request) []ControllerHandler {
	uri := req.URL.Path
	method := req.Method
	upperMethod := strings.ToUpper(method)
	if methodHandlers, ok := c.router[upperMethod]; ok {
		return methodHandlers.FindHandler(uri)
	}
	return nil
}

// Use 增加中间件
func (c *Core) Use(middlewares ...ControllerHandler) {
	c.middleswares = append(c.middleswares, middlewares...)
}
