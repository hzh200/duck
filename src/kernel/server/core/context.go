package core

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Response http.ResponseWriter
	Request *http.Request
	RouteParams map[string]string
}

func (context *Context) GetParam(key string) string {
	return context.Request.URL.Query().Get(key)
}

func (context *Context) PostFormParam(key string) string {
	context.Request.ParseForm()
	return context.Request.FormValue(key)
}

func (context *Context) BodyParams() map[string]interface{} {
	m := make(map[string]interface{})
	json.NewDecoder(context.Request.Body).Decode(&m)
	return m
}

func (context *Context) Status(code int) {
	context.Response.WriteHeader(code)
}

func (context *Context) SetHeader(key string, value string) {
	context.Response.Header().Add(key, value)
}

func (context *Context) Raw(data []byte) {
	context.Status(http.StatusOK)
	context.Response.Write(data)
}

func (context *Context) Text(text string) {
	context.Status(http.StatusOK)
	context.SetHeader("Content-Type", "text/plain")
	context.Response.Write([]byte(text))
}

func (context *Context) HTML(html string) {
	context.Status(http.StatusOK)
	context.SetHeader("Content-Type", "text/html")
	context.Response.Write([]byte(html))
}

func (context *Context) Json(data J) error {
	context.Status(http.StatusOK)
	context.SetHeader("Content-Type", "application/json")
	jsonData, err := json.Marshal(data) 
	if err != nil {
		return err
	}
	context.Response.Write(jsonData)
	return nil
}
