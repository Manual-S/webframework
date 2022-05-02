// 框架提供的中间件
package middleware

import (
	"log"
	"time"
	"webframework/framework"
)

func Cost() framework.ControllerHandler {
	return func(c *framework.Context) error {
		start := time.Now()

		c.Next()

		end := time.Now()
		cost := end.Sub(start)

		log.Printf("api uri :%v cost :%v", c.GetRequest().URL, cost.Microseconds())
		return nil
	}
}
