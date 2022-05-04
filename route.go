// 路由层
package main

import (
	"webframework/framework"
)

func registerRouter(core *framework.Core) {
	//core.Use(middleware.Cost())

	// 需求1和2 http方法和静态路由匹配
	core.Get("/user/login", UserLoginController)
	core.Get("/foo", FooControllerHandler)
	core.Post("/user/login", UserLoginController)

	// 需求3 批量通用前缀
	subjectApi := core.Group("/subject")
	{
		// 需求4:动态路由
		//subjectApi.Delete("/:id", SubjectDelController)
		//subjectApi.Put("/:id", SubjectUpdateController)
		subjectApi.Get("/:id", SubjectGetController)
		//subjectApi.Get("/list/all", SubjectListController)
	}
}
