// 框架提供的中间件
package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"webframework/framework"
)

func Recovery() framework.ControllerHandler {
	return func(c *framework.Context) error {
		defer func() {
			if err := recover(); err != nil {
				c.SetStatus(http.StatusInternalServerError).Json(err)
			}
		}()

		c.Next()

		return nil
	}
}

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

// TimeHandler 这种实现中间件的方式是通过函数嵌套
func TimeHandler(fun framework.ControllerHandler, d time.Duration) framework.ControllerHandler {
	// 返回一个匿名函数
	return func(c *framework.Context) error {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d*time.Second)
		defer cancel()

		// 创建一个goroutine来处理具体的业务逻辑
		go func() {

			defer func() {
				// 异常处理
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			fun(c)

			finish <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			// 发生了异常
			c.WriteMux().Lock()
			defer c.WriteMux().Unlock()
			log.Println(p)
			c.SetStatus(http.StatusInternalServerError).Json("panic")
		case <-finish:
			fmt.Println("finish")
		case <-durationCtx.Done():
			c.WriteMux().Lock()
			defer c.WriteMux().Unlock()
			c.SetStatus(http.StatusInternalServerError).Json("time out")
			c.SetHasTimeout()
		}

		return nil
	}
}
