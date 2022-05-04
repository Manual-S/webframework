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

// ServeHTTP 实现ServeHTTP接口 就可以将所有的http请求转向我们自己的处理逻辑
// 在ServeHTTP中 可以自由定义路由映射规则
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request, response)

	// 寻找路由
	router := c.FindRouteByRequest(request)
	if router == nil {
		ctx.SetOkStatus().Json("not found")
		return
	}

	// 设置context中的handlers字段
	ctx.SetHandlers(router.handler)

	// 设置路由参数
	params := router.parseParmsFromEndNode(request.URL.Path)
	ctx.SetParams(params)

	err := ctx.Next()
	if err != nil {
		ctx.SetStatus(http.StatusInternalServerError).Json("inner error")
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

func (c *Core) FindRouteByRequest(req *http.Request) *node {
	uri := req.URL.Path
	method := req.Method
	upperMethod := strings.ToUpper(method)
	if methodHandlers, ok := c.router[upperMethod]; ok {
		return methodHandlers.root.matchNode(uri)
	}
	return nil
}

// Use 增加中间件
func (c *Core) Use(middlewares ...ControllerHandler) {
	c.middleswares = append(c.middleswares, middlewares...)
}

func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}
