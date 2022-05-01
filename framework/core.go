package framework

import (
	"log"
	"net/http"
	"strings"
)

type Core struct {
	router map[string]*Tree
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

	err := router(ctx)
	if err != nil {
		ctx.Json(http.StatusInternalServerError, "500")
		return
	}
}

func (c *Core) Get(url string, handler ControllerHandler) {
	err := c.router["GET"].AddRouter(url, handler)
	if err != nil {
		log.Fatal("add router error" + url)
	}
}

func (c *Core) Post(url string, handler ControllerHandler) {

}

func (c *Core) FindRouteByRequest(req *http.Request) ControllerHandler {
	uri := req.URL.Path
	method := req.Method
	upperMethod := strings.ToUpper(method)
	if methodHandlers, ok := c.router[upperMethod]; ok {
		return methodHandlers.FindHandler(uri)
	}
	return nil
}
