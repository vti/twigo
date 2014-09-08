package action

import (
	"net/http"

	"github.com/vti/twigo/app/model"
)

type ViewTag struct {
	BaseAction
}

func (action *ViewTag) Execute(w http.ResponseWriter, r *http.Request) {
	home := action.Context.App.Home

	tag := action.Context.Capture("tag")

	dm := &model.DocumentManager{Root: home + "/articles/"}
	documents, err := dm.LoadDocumentsByTag(tag)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	action.Context.SetTemplateName("layouts/html")
	action.Context.SetTemplateFiles([]string{
		"layouts/html.tpl",
		"tag.tpl",
		"article-meta.tpl",
	})

	vars := map[string]interface{}{
		"Conf":      action.Context.App.Conf,
		"Tag":       tag,
		"Documents": documents,
	}
	action.Context.SetTemplateVars(vars)
}
