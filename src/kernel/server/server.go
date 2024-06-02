package server

import (
	"duck/kernel/server/core"
)

func StartServer(port uint16) error {
	engine := core.New(port)
	
	engine.Get("/", (func(context *core.Context) {
		context.HTML("<strong>hello</strong>")
	}))
	engine.Get("/task/:taskNo", (func(context *core.Context) {
		
	}))
	engine.Post("/task", (func(context *core.Context) {
		context.PostParam("filename")
	}))
	engine.Post("/tasks", (func(context *core.Context) {
		context.PostParam("groupname")
	}))
	
	// return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return engine.Run()
}
