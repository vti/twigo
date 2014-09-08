package action

import (
	"fmt"
	"net/http"

	"github.com/gorilla/feeds"
	"github.com/vti/twigo/app/model"
)

type ListArticlesByTagRss struct {
	BaseAction
}

func (action *ListArticlesByTagRss) Execute(w http.ResponseWriter, r *http.Request) {
	home := action.Context.App.Home

	tag := action.Context.Capture("tag")

	dm := &model.DocumentManager{Root: home + "/articles/"}
	documents, err := dm.LoadDocumentsByTag(tag)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	pubDate := model.Date{}
	if len(documents) > 0 {
		pubDate = documents[0].Created
	}

	conf := action.Context.App.Conf
	feed := &feeds.Feed{
		Title:       conf.Title,
		Link:        &feeds.Link{Href: action.buildUrl(r, "Index")},
		Description: conf.Description,
		Author:      &feeds.Author{conf.Author, ""},
		Created:     pubDate.Time(),
	}

	for _, document := range documents {
		description := document.Preview
		if description == "" {
			description = document.Content
		}

		feed.Items = append(feed.Items,
			&feeds.Item{
				Title:       document.Meta["Title"],
				Link:        &feeds.Link{Href: action.buildViewArticleUrl(r, document)},
				Id:          action.buildViewArticleUrl(r, document),
				Description: description,
				Author:      &feeds.Author{conf.Author, ""},
				Created:     document.Created.Time(),
			})
	}

	rss, err := feed.ToRss()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/rss+xml")
	fmt.Fprintf(w, rss)
}

func (action *ListArticlesByTagRss) buildUrl(r *http.Request, name string, args ...string) string {
	route := action.Context.App.Router.Get(name)
	if route == nil {
		return ""
	}
	url, err := route.URL(args...)
	if err != nil {
		return ""
	}
	return "http://" + r.Host + url.String()
}

func (action *ListArticlesByTagRss) buildViewArticleUrl(r *http.Request, document *model.Document) string {
	route := action.Context.App.Router.Get("ViewArticle")
	if route == nil {
		return ""
	}
	url, err := route.URL("year", document.Created["Year"].String(),
		"month", document.Created["Month"].String(),
		"title", document.Slug)
	if err != nil {
		return ""
	}
	return "http://" + r.Host + url.String()
}
