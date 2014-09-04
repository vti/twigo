package app

import (
    "github.com/gorilla/mux"
)

type Configuration struct {
	Title       string
	Author      string
	Generator   string
	Description string
	About       string
}

type Twigo struct {
	Home string
	Conf *Configuration
    Router *mux.Router
}
