package utils

import (
	"github.com/gorilla/mux"

	"github.com/vti/twigo/app/model"
)

func BuildViewArticleUrl(router *mux.Router, document *model.Document) string {
	return BuildUrl(router, "ViewArticle", "year", document.Created["Year"].String(),
		"month", document.Created["Month"].String(),
		"title", document.Slug)
}

func BuildUrl(router *mux.Router, name string, args ...string) string {
	route := router.Get(name)
	if route == nil {
		return ""
	}
	url, err := route.URL(args...)
	if err != nil {
		return ""
	}
	return url.String()
}
