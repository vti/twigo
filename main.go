package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/gorilla/mux"

	"github.com/vti/twigo/app"
	"github.com/vti/twigo/app/action"
	"github.com/vti/twigo/app/model"
	"github.com/vti/twigo/app/utils"
)

func detectHome() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func main() {
	usage := `twigo.

Usage:
  twigo serve --conf=<conf> --listen=<listen>
  twigo -h | --help

Options:
  -h --help         Show this screen.
  --conf=<conf>     Path to configuration file (conf.json).
  --listen=<listen> Listen options (:8080).`

	arguments, err := docopt.Parse(usage, nil, true, "twigo", false)
	if err != nil {
		log.Fatal("error parsing command-line options:", err)
	}

	conf := utils.LoadConfiguration(arguments["--conf"].(string))
	home := detectHome()
	router := mux.NewRouter()

	twigo := &app.Twigo{Conf: conf, Home: home, Router: router}

	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	http.Handle("/static/", fileServer)

	router.Handle("/articles/{year:[0-9]{4}}/{month:0?[1-9]|1[012]}/{title:[A-Za-z0-9-]+}.html",
		makeHandler(&action.ViewArticle{}, twigo)).
		Methods("GET").Name("ViewArticle")
	router.Handle("/pages/{title:[a-z0-9-]+}.html",
		makeHandler(&action.ViewPage{}, twigo)).
		Methods("GET")
	router.Handle("/tags.html",
		makeHandler(&action.ListTags{}, twigo)).
		Methods("GET")
	router.Handle("/tags/{tag:[A-Za-z0-9]+}.html",
		makeHandler(&action.ListArticlesByTag{}, twigo)).
		Methods("GET").Name("ListArticlesByTag")
	router.Handle("/tags/{tag:[A-Za-z0-9]+}.rss",
		makeHandler(&action.ListArticlesByTagRss{}, twigo)).
		Methods("GET").Name("ListArticlesByTagRss")
	router.Handle("/",
		makeHandler(&action.ListArticles{}, twigo)).
		Methods("GET").Name("Index")
	router.Handle("/index.rss",
		makeHandler(&action.ListArticlesRss{}, twigo)).
		Methods("GET").Name("ListArticlesRss")
	router.Handle("/archive.html",
		makeHandler(&action.ListArticlesArchive{}, twigo)).
		Methods("GET").Name("ListArticlesArchive")

	http.Handle("/", router)

	http.ListenAndServe(arguments["--listen"].(string), nil)
}

func makeHandler(action app.Action, a *app.Twigo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		context := app.NewContext(a, r)

		action.SetContext(context)
		action.Execute(w, r)

		if context.TemplateName != "" {
			var templateFiles []string
			for _, file := range context.TemplateFiles {
				templateFiles = append(templateFiles, context.App.Home+"/templates/"+file)
			}

			t, err := template.New(context.TemplateName).
				Funcs(template.FuncMap{
				"safeHtml": func(text string) template.HTML {
					return template.HTML(text)
				},
				"conf": func() *app.Configuration {
					return context.App.Conf
				},
				"partial": func(name string, ctx interface{}) template.HTML {
					name = context.App.Home + "/templates/" + name
					t, err := template.New(name).ParseFiles(name)
					if err != nil {
						return ""
					}
					output := bytes.Buffer{}
					err = t.Execute(&output, ctx)
					if err != nil {
						return ""
					}
					fmt.Println(output)
					return template.HTML(output.String())
				},
				"dateFmt": func(date model.Date) string {
					const layout = "Mon, 2 Jan 2006"
					t := time.Date(int(date["Year"]), time.Month(date["Month"]), int(date["Day"]), 0, 0, 0, 0, time.Local)
					return t.Format(layout)
				},
				"buildViewArticleUrl": func(document *model.Document) string {
					return "http://" + r.Host + utils.BuildViewArticleUrl(context.App.Router, document)
				},
				"buildUrl": func(name string, args ...string) string {
					return "http://" + r.Host + utils.BuildUrl(context.App.Router, name, args...)
				}}).
				ParseFiles(templateFiles...)

			if err != nil {
				log.Print(err)
				http.NotFound(w, r)
				return
			}

			err = t.Execute(w, context.TemplateVars)
			if err != nil {
				log.Print(err)
			}
		}
	}
}
