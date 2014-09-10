package utils

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseJSON(t *testing.T) {
	var json = `
{
    "title":"foo",
    "page_limit": 10
}
`

	conf, _ := parseJSON([]byte(json))

	assert.Equal(t, "foo", conf.Title)
	assert.Equal(t, 10, conf.PageLimit)
}

func TestParseYAML(t *testing.T) {
	var yaml = `
---
title: foo
page_limit: 10
`

	conf, _ := parseYAML([]byte(yaml))

	assert.Equal(t, "foo", conf.Title)
	assert.Equal(t, 10, conf.PageLimit)
}

func TestSlurp(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	ioutil.WriteFile(tempDir+"config.yaml", []byte(string("foobarbaz")), 0644)

	conf, _ := slurp(tempDir + "config.yaml")

	assert.Equal(t, "foobarbaz", string(conf))

	os.RemoveAll(tempDir)
}

func TestLoadConfigurationFromYAML(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	var yaml = `
---
title: foo
page_limit: 10
`

	ioutil.WriteFile(tempDir+"config.yaml", []byte(yaml), 0644)

	conf := LoadConfiguration(tempDir + "config.yaml")

	assert.Equal(t, "foo", conf.Title)
	assert.Equal(t, 10, conf.PageLimit)

	os.RemoveAll(tempDir)
}

func TestLoadConfigurationFromJSON(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	var yaml = `
{
"title": "foo",
"page_limit": 10
}
`

	ioutil.WriteFile(tempDir+"config.json", []byte(yaml), 0644)

	conf := LoadConfiguration(tempDir + "config.json")

	assert.Equal(t, "foo", conf.Title)
	assert.Equal(t, 10, conf.PageLimit)

	os.RemoveAll(tempDir)
}
