package core

import (
	"net/http"
	"strings"
	"testing"
)

func TestRouterAddRoute(t *testing.T) {
	router := NewRouter()

	for _, routes := range router.GetRoutes() {
		if len(routes) != 0 {
			t.Fatalf("router routers are not empty right after being created")
		}
	}

	path := "/hello"
	var handler Handler = func (context *Context) {}
	err := router.AddRoute(http.MethodGet, path, &handler)
	if err != nil {
		t.Fatalf("router add route failed: %s", err)
	}
	storedRoutes := router.GetRoutes()
	if len(storedRoutes[http.MethodGet]) != 1 {
		t.Fatalf("router didn't store added path properly, length of storedRoutes: %d",len(storedRoutes) )
	}
	if strings.Compare(storedRoutes[http.MethodGet][0], path) != 0 {
		t.Fatalf("router didn't store added path properly, stored %s, get %s", path, storedRoutes[http.MethodGet][0])
	}
}

func TestRouterRoute(t *testing.T) {
	router := NewRouter()
	path := "/hello"
	var handler Handler = func (context *Context) {}
	router.AddRoute(http.MethodGet, path, &handler)
	if found, routeHandler := router.Route(http.MethodGet, path); !found || routeHandler != &handler {
		t.Fatalf("router route failed, %p", routeHandler)
	}
}
