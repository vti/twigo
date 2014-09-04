package action

import (
	"html/template"
	"net/http"

	"github.com/vti/twigo/app"
	"github.com/vti/twigo/app/model"
)

type ViewPage struct {
	Context *app.Context
}

func (action *ViewPage) SetContext(context *app.Context) {
	action.Context = context
}

func (action *ViewPage) Execute(w http.ResponseWriter, r *http.Request) {
	home := action.Context.App.Home
	title := action.Context.Capture("title")

	dm := &model.DocumentManager{Root: home + "/pages/"}
	document, err := dm.LoadDocument(title)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	action.Context.SetTemplateName("layouts/html")
	action.Context.SetTemplateFiles([]string{"layouts/html.tpl", "page.tpl"})

	vars := map[string]interface{}{
		"Conf": action.Context.App.Conf,
		"Document": map[string]interface{}{
			"Meta":    document.Meta,
			"Content": template.HTML(document.Content),
		},
	}

	action.Context.SetTemplateVars(vars)
}