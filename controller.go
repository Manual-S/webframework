package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"webframework/framework"
)

func FooControllerHandler(c *framework.Context) error {

	finish := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)

	durationCtx, cancel := context.WithTimeout(c.BaseContext(), 1*time.Second)
	defer cancel()

	// 创建一个goroutine来处理具体的业务逻辑
	go func() {

		defer func() {
			// 异常处理
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()

		time.Sleep(10 * time.Second)
		c.Json(http.StatusOK, "ok")

		finish <- struct{}{}
	}()

	select {
	case p := <-panicChan:
		// 发生了异常
		log.Println(p)
		c.Json(http.StatusInternalServerError, "panic")
	case <-finish:
		fmt.Println("finish")
	case <-durationCtx.Done():
		c.Json(http.StatusInternalServerError, "time out")
	}

	return nil
}
