package action

import (
	"html/template"
	"log"
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

	t, err := template.New("layouts/html").ParseFiles(home+"/templates/layouts/html.tpl",
		home+"/templates/page.tpl")
	if err != nil {
		log.Print(err)
		http.NotFound(w, r)
		return
	}

	vars := map[string]interface{}{
		"Conf": action.Context.App.Conf,
		"Document": map[string]interface{}{
			"Meta":    document.Meta,
			"Content": template.HTML(document.Content),
		},
	}

	t.Execute(w, vars)
}
