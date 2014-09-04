package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Context struct {
	App *Twigo
	Req *http.Request
}

func NewContext(app *Twigo, r *http.Request) *Context {
	return &Context{App: app, Req: r}
}

func (context *Context) Capture(name string) string {
	return mux.Vars(context.Req)[name]
}
