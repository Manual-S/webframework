// 路由层
package main

import "webframework/framework"

func registerRouter(core *framework.Core) {
	core.Get("foo", FooControllerHandler)
}
