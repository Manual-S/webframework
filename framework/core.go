package framework

import "net/http"

type Core struct {
	router map[string]ControllerHandler
}

func NewCore() *Core {
	return &Core{
		router: map[string]ControllerHandler{},
	}
}

func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request, response)
	// 路由选择
	router := c.router["foo"]
	// 强制类型转换
	router(ctx)
}

func (c *Core) Get(url string, handler ControllerHandler) {
	c.router[url] = handler
}
