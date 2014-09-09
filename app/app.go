package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Configuration struct {
	Title       string
	Author      string
	Generator   string
	Description string
	About       string
	Menu        []map[string]string
	PageLimit   int
	Footer      string
	BaseUrl     string
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
