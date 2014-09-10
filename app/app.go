package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

type DisqusType struct {
	Shortname string `json:"shortname"       yaml:"shortname"`
	Developer int    `json:"developer"       yaml:"developer"`
}

type Configuration struct {
	Title       string              `json:"title"       yaml:"title"`
	Author      string              `json:"author"      yaml:"author"`
	Generator   string              `json:"generator"   yaml:"generator"`
	Description string              `json:"description" yaml:"description"`
	About       string              `json:"about"       yaml:"about"`
	Menu        []map[string]string `json:"menu"        yaml:"menu"`
	PageLimit   int                 `json:"page_limit"  yaml:"page_limit"`
	Footer      string              `json:"footer"      yaml:"footer"`
	BaseUrl     string              `json:"base_url"    yaml:"base_url"`
	Disqus      DisqusType
}

type Twigo struct {
	Home   string
	Conf   *Configuration
	Router *mux.Router
}

type Action interface {
	SetContext(*Context)
	Execute(w http.ResponseWriter, r *http.Request)
}
