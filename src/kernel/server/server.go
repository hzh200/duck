package server

import (
	"duck/kernel/server/core"
)

func Run(port uint16) error {
	server := core.New(port)
	
	server.Get("/", (func(context *core.Context) {
		context.HTML("<strong>hello</strong>")
	}))
	server.Get("/task/:taskNo", (func(context *core.Context) {
		
	}))
	server.Post("/task", (func(context *core.Context) {
		context.PostParam("filename")
	}))
	server.Post("/tasks", (func(context *core.Context) {
		context.PostParam("groupname")
	}))
	
	// return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return server.Start()
}
