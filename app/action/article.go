package action

import (
	"html/template"
	"net/http"

	"github.com/vti/twigo/app"
	"github.com/vti/twigo/app/model"
)

type ViewArticle struct {
	Context *app.Context
}

func (action *ViewArticle) SetContext(context *app.Context) {
	action.Context = context
}

func (action *ViewArticle) Execute(w http.ResponseWriter, r *http.Request) {
	home := action.Context.App.Home

	title := action.Context.Capture("title")
	year := action.Context.Capture("year")
	month := action.Context.Capture("month")

	dm := &model.DocumentManager{Root: home + "/articles/"}
	document, err := dm.FindDocument(title, year, month)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if document == nil {
		http.NotFound(w, r)
		return
	}

	action.Context.SetTemplateName("layouts/html")
	action.Context.SetTemplateFiles([]string{
		"layouts/html.tpl",
		"article.tpl",
		"article-meta.tpl",
	})

	vars := map[string]interface{}{
		"Conf": action.Context.App.Conf,
		"Document": map[string]interface{}{
			"Meta":    document.Meta,
			"Preview": template.HTML(document.Preview),
			"Content": template.HTML(document.Content),
		},
	}
	action.Context.SetTemplateVars(vars)
}
