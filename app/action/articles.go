package action

import (
	"net/http"

	"github.com/vti/twigo/app/model"
)

type ListArticles struct {
	BaseAction
}

func (action *ListArticles) Execute(w http.ResponseWriter, r *http.Request) {
	home := action.Context.App.Home

	limit := action.Context.App.Conf.PageLimit
	if limit == 0 {
		limit = 10
	}
	offset := r.URL.Query().Get("timestamp")

	dm := &model.DocumentManager{Root: home + "/articles/"}
	documents, err := dm.LoadDocuments(limit, offset)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	action.Context.SetTemplateName("layouts/html")
	action.Context.SetTemplateFiles([]string{
		"layouts/html.tpl",
		"articles.tpl",
		"article-meta.tpl",
	})

	vars := map[string]interface{}{
		"Conf":           action.Context.App.Conf,
		"Documents":      documents,
		"PrevPageOffset": dm.PrevPageOffset(limit, offset),
		"NextPageOffset": dm.NextPageOffset(limit, offset),
	}
	action.Context.SetTemplateVars(vars)
}
