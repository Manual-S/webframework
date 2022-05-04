// 理解成server层
// 当前框架的一些demo
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
		c.SetOkStatus()

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

type User struct {
	UserName string `json:"user_name"`
	Pasword  string `json:"pasword"`
}

func UserLoginController(c *framework.Context) error {
	user := User{}
	c.BindJson(&user)
	c.SetOkStatus().Json(user)
	return nil
}

func SubjectGetController(c *framework.Context) error {
	// 具体的业务逻辑
	hash := make(map[string]int)
	id, _ := c.ParamInt("id", 0)
	hash["id"] = id
	c.SetOkStatus().Json(hash)

	return nil
}

func FooControllerHandler2(c *framework.Context) error {
	time.Sleep(10 * time.Second)
	c.SetOkStatus().Json("ok")
	return nil
}
