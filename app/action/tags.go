package action

import (
	"net/http"

	"github.com/vti/twigo/app/model"
)

type ListTags struct {
	BaseAction
}

func (action *ListTags) Execute(w http.ResponseWriter, r *http.Request) {
	home := action.Context.App.Home

	dm := &model.DocumentManager{Root: home + "/articles/"}
	documents, err := dm.LoadDocuments(0, "")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	tags := map[string]int{}
	for _, document := range documents {
		for _, t := range document.Tags {
			tags[t]++
		}
	}

	action.Context.SetTemplateName("layouts/html")
	action.Context.SetTemplateFiles([]string{
		"layouts/html.tpl",
		"tags.tpl",
	})

	vars := map[string]interface{}{
		"Conf": action.Context.App.Conf,
		"Tags": tags,
	}
	action.Context.SetTemplateVars(vars)
}
