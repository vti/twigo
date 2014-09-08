package action

import (
	"fmt"
	"net/http"

	"github.com/gorilla/feeds"

	"github.com/vti/twigo/app/model"
	"github.com/vti/twigo/app/utils"
)

type ListArticlesRss struct {
	BaseAction
}

func (action *ListArticlesRss) Execute(w http.ResponseWriter, r *http.Request) {
	home := action.Context.App.Home
	router := action.Context.App.Router

	dm := &model.DocumentManager{Root: home + "/articles/"}
	documents, err := dm.LoadDocuments(action.Context.App.Conf.PageLimit, "")
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
		Link:        &feeds.Link{Href: utils.BuildUrl(router, "Index")},
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
				Link:        &feeds.Link{Href: utils.BuildViewArticleUrl(router, document)},
				Id:          utils.BuildViewArticleUrl(router, document),
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
