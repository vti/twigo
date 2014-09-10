package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v1"

	"github.com/vti/twigo/app"
)

func LoadConfiguration(path string) *app.Configuration {
	ext := filepath.Ext(path)

	var parser func(data []byte) (*app.Configuration, error)
	if ext == ".json" {
		parser = parseJSON
	} else if ext == ".yml" || ext == ".yaml" {
		parser = parseYAML
	} else {
		log.Fatal("Don't know how to load config file with ext=", ext)
	}

	data, err := slurp(path)
	if err != nil {
		log.Fatal(err)
	}

	conf, err := parser(data)
	if err != nil {
		log.Fatal(err)
	}

	return conf
}

func slurp(path string) (data []byte, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err = ioutil.ReadAll(file)
	return
}

func parseJSON(data []byte) (*app.Configuration, error) {
	conf := app.Configuration{}

	err := json.Unmarshal(data, &conf)

	return &conf, err
}

func parseYAML(data []byte) (*app.Configuration, error) {
	conf := app.Configuration{}

	err := yaml.Unmarshal(data, &conf)

	return &conf, err
}
