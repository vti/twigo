package utils

import (
	"encoding/json"
	"log"
	"os"

	"github.com/vti/twigo/app"
)

func LoadConfiguration(path string) *app.Configuration {
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
