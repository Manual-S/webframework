package main

import (
	"net/http"
	"webframework/framework"
)

func main() {
	core := framework.NewCore()
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
