package main

import (
	"bytes"
	"encoding/json"
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
)

func detectHome() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func loadConfiguration(path string) *app.Configuration {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := app.Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err)
	}

	return &configuration
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

	conf := loadConfiguration(arguments["--conf"].(string))
	home := detectHome()
	router := mux.NewRouter()

	twigo := &app.Twigo{Conf: conf, Home: home, Router: router}

	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	http.Handle("/static/", fileServer)

	router.Handle("/articles/{year:[0-9]{4}}/{month:0?[1-9]|1[012]}/{title:[A-Za-z0-9-]+}",
		makeHandler(&action.ViewArticle{}, twigo)).
		Methods("GET").Name("ViewArticle")
	router.Handle("/pages/{title:[a-z0-9]+}",
		makeHandler(&action.ViewPage{}, twigo)).
		Methods("GET")
	router.Handle("/tags",
		makeHandler(&action.ListTags{}, twigo)).
		Methods("GET")
	//router.Handle("/tags/{tag:[a-z0-9]+}",
		//makeHandler(&action.ViewTag{}, twigo)).
		//Methods("GET")
	router.Handle("/",
		makeHandler(&action.ListArticles{}, twigo)).
		Methods("GET")

	http.Handle("/", router)

	http.ListenAndServe(arguments["--listen"].(string), nil)
}

type Action interface {
	SetContext(*app.Context)
	Execute(w http.ResponseWriter, r *http.Request)
}

func makeHandler(action Action, a *app.Twigo) http.HandlerFunc {
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
				"buildViewArticleUrl": func(document *model.Document) string {
					route := context.App.Router.Get("ViewArticle")
					if route == nil {
						return ""
					}
					url, err := route.URL("year", document.Created["Year"].String(),
						"month", document.Created["Month"].String(),
						"title", document.Slug)
					if err != nil {
						return ""
					}
					return url.String()
				},
				"dateFmt": func(date model.Date) string {
					const layout = "Mo, 2 Jan 2006"
					t := time.Date(int(date["Year"]), time.Month(date["Month"]), int(date["Day"]), 0, 0, 0, 0, time.Local)
					return t.Format(layout)
				},
				"buildUrl": func(name string, args ...string) string {
					route := context.App.Router.Get(name)
					if route == nil {
						return ""
					}
					url, err := route.URL(args...)
					if err != nil {
						return ""
					}
					return url.String()
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
