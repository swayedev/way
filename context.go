package way

import (
	"net/http"
)

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{Response: w, Request: r}
}
