package core

import (
	"errors"
	"net/http"
	"strings"
)

type Node struct {
	val string
	children []*Node
	isLeaf bool
	handler *Handler
}

type Router struct {
	root *Node
}

func NewRouter() *Router {
	root := Node{
		val: "",
		children: make([]*Node, 0),
		isLeaf: false,
		handler: nil,
	}
	router := Router{&root}
	return &router
}

func (router *Router) AddRoute(method string, path string, handler *Handler) error {
	if len(path) < 1 || path[0] != '/' || strings.Contains(path, "//") {
		return errors.New("invalid path")
	}
	if len(path) > 1 && path[len(path) - 1] == '/' {
		path = path[:len(path) - 1]
	}
	curNode := router.root
	subPaths := strings.Split(method + path, "/")
	for i, subPath := range subPaths {
		before := curNode
		for _, child := range curNode.children {
			if strings.Compare(child.val, subPath) == 0 {
				curNode = child
				break
			}
		}
		if before == curNode {
			newNode := Node{
				val: subPath,
				children: make([]*Node, 0),
				isLeaf: false,
				handler: nil,
			}
			curNode.children = append(curNode.children, &newNode)
			curNode = &newNode
		}
		if i == len(subPaths) - 1 {
			curNode.isLeaf = true
			curNode.handler = handler
		}
	}
	return nil
}

func (router *Router) Route(method string, path string) (bool, *Handler) {
	if len(path) < 1 || path[0] != '/' || strings.Contains(path, "//") {
		return false, nil
	}

	if len(path) > 1 && path[len(path) - 1] == '/' {
		path = path[:len(path) - 1]
	}
	curNode := router.root
	subPaths := strings.Split(method + path, "/")
	for i, subPath := range subPaths {
		before := curNode
		for _, child := range curNode.children {
			if strings.Compare(child.val, subPath) == 0 {
				curNode = child
				break
			}
		}
		if before == curNode {
			return false, nil
		}
		if curNode.isLeaf && i == len(subPaths) - 1 {
			return true, curNode.handler
		}
	}
	return false, nil
}

func (router *Router) GetRoutes() map[string][]string {
	routes := make(map[string][]string)
	routes[http.MethodGet] = make([]string, 0)
	routes[http.MethodPost] = make([]string, 0)
	router.getRoutesCore(router.root, make([]string, 0), routes)
	return routes
}

func (router *Router) getRoutesCore(curNode *Node, curRoute []string, routes map[string][]string) {
	if curNode.isLeaf && len(curNode.children) == 0 {
		routes[curRoute[0]] = append(routes[curRoute[0]], "/" + strings.Join(curRoute[1:], "/"))		
	}
	for _, child := range curNode.children {
		router.getRoutesCore(child, append(curRoute, child.val), routes)
	}
}
