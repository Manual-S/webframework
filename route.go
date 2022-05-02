// 路由层
package main

import (
	"webframework/framework"
	"webframework/framework/middleware"
)

func registerRouter(core *framework.Core) {
	core.Use(middleware.Cost())
	core.Get("/user/login", UserLoginController)
	//subjectApi := core.Group("/subject")
	//{
	//	// 需求4:动态路由
	//	subjectApi.Delete("/:id", SubjectDelController)
	//	subjectApi.Put("/:id", SubjectUpdateController)
	//	subjectApi.Get("/:id", SubjectGetController)
	//	subjectApi.Get("/list/all", SubjectListController)
	//}
}
