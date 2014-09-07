package action

import (
	"net/http"

	"github.com/vti/twigo/app"
	"github.com/vti/twigo/app/model"
)

type ViewTag struct {
	Context *app.Context
}

func (action *ViewTag) SetContext(context *app.Context) {
	action.Context = context
}

func (action *ViewTag) Execute(w http.ResponseWriter, r *http.Request) {
	home := action.Context.App.Home

	tag := action.Context.Capture("tag")

	dm := &model.DocumentManager{Root: home + "/articles/"}
	documents, err := dm.LoadDocuments(0, "")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	taggedDocuments := []*model.Document{}
	for _, document := range documents {
		for _, t := range document.Tags {
			if t == tag {
				taggedDocuments = append(taggedDocuments, document)
				break
			}
		}
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
		"Documents": taggedDocuments,
	}
	action.Context.SetTemplateVars(vars)
}
