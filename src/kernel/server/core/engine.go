package core

import (
	"fmt"
	"net/http"
)

type Engine struct {
	router *Router
	port uint16
}

func New(port uint16) *Engine {
	engine := Engine{}
	engine.router = NewRouter()
	engine.port = port
	return &engine
}

func (engine *Engine) addRoute(method string, url string, handler Handler) {
	engine.router.AddRoute(method, url, &handler)
}

func (engine *Engine) Get(url string, handler Handler) {
	engine.addRoute(http.MethodGet, url, handler)
}

func (engine *Engine) Post(url string, handler Handler) {
	engine.addRoute(http.MethodPost, url, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	found, handler := engine.router.Route(r.Method, r.URL.Path)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	context := Context{}
	context.Response = w
	context.Request = r
	(*handler)(&context)
}

func (engine *Engine) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", int(engine.port)), engine)
}
