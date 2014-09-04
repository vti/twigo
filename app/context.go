package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Context struct {
	App           *Twigo
	Req           *http.Request
	TemplateName  string
	TemplateFiles []string
	TemplateVars  map[string]interface{}
}

func NewContext(app *Twigo, r *http.Request) *Context {
	return &Context{App: app, Req: r}
}

func (context *Context) Capture(name string) string {
	return mux.Vars(context.Req)[name]
}

func (context *Context) SetTemplateName(name string) *Context {
	context.TemplateName = name
	return context
}

func (context *Context) SetTemplateFiles(files []string) *Context {
	context.TemplateFiles = files
	return context
}

func (context *Context) SetTemplateVars(vars map[string]interface{}) *Context {
	context.TemplateVars = vars
	return context
}
