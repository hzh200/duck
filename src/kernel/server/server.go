package server

import (
	"duck/kernel/server/core"
)

func StartServer(port uint16) error {
	engine := core.New(port)

	for route, handler := range GetRoutes {
		engine.Get(route, handler)
	}

	for route, handler := range PostRoutes {
		engine.Post(route, handler)
	}
	
	// return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return engine.Run()
}
