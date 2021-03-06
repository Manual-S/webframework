package framework

type IGroup interface {
	Get(string, ControllerHandler)
	Post(string, ControllerHandler)
	//Put(string, ControllerHandler)
	//Delete(string, ControllerHandler)
}

type Group struct {
	core   *Core
	prefix string
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:   core,
		prefix: prefix,
	}
}

func (g *Group) Get(uri string, handler ControllerHandler) {
	uri = g.prefix + uri
	g.core.Get(uri, handler)
}

func (g *Group) Post(uri string, handler ControllerHandler) {
	uri = g.prefix + uri
	g.core.Post(uri, handler)
}
