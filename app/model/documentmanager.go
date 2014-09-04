package model

import (
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

type Document struct {
	Meta    map[string]string
	Preview string
	Content string
}

type DocumentManager struct {
	Root string
}

func (dm *DocumentManager) LoadDocument(name string) (*Document, error) {
	path := dm.Root + name + ".markdown"

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	input, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return parseInput(input)
}

func (dm *DocumentManager) FindDocument(name string, year string, month string) (*Document, error) {
	files, err := ioutil.ReadDir(dm.Root)
	if err != nil {
		return nil, err
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return nil, err
	}
	monthInt, err := strconv.Atoi(month)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		if !strings.HasSuffix(f.Name(), name+".markdown") {
			continue
		}

		date, fileName := parseFileName(f.Name())
		fileName = strings.TrimSuffix(fileName, ".markdown")

		if fileName != name {
			continue
		}

		if yearInt > 0 && monthInt > 0 {
			if date == nil {
				continue
			} else if yearInt != date["year"] || monthInt != date["month"] {
				continue
			}
		}

		fullName := strings.TrimSuffix(f.Name(), ".markdown")
		return dm.LoadDocument(fullName)
	}

	return nil, nil
}

func (dm *DocumentManager) ListDocuments() ([]*Document, error) {
	files, err := ioutil.ReadDir(dm.Root)
	if err != nil {
		return nil, err
	}

    var documents []*Document
	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		if !strings.HasSuffix(f.Name(), ".markdown") {
			continue
		}

		fullName := strings.TrimSuffix(f.Name(), ".markdown")
        document, err := dm.LoadDocument(fullName)
        if err != nil {
            continue
        }

		documents = append(documents, document)
	}

	return documents, nil
}

func parseFileName(input string) (map[string]int, string) {
	dateRe := regexp.MustCompile("^[0-9]{8}-")

	if dateRe.MatchString(input) {
		parts := strings.SplitN(input, "-", 2)

		date, _ := parseDate(parts[0])

		if date != nil {
			return date, parts[1]
		}
	}

	return nil, input
}

func parseDate(input string) (map[string]int, error) {
	t, err := time.Parse("20060102", input)
	if err != nil {
		return nil, err
	}

	date := map[string]int{"year": t.Year(),
		"month": int(t.Month()),
		"day":   t.Day()}
	return date, nil
}

func parseInput(input []byte) (*Document, error) {
	metaRaw, contentRaw := splitMetaAndContent(input)
	meta := parseMeta(metaRaw)
	preview, content := splitPreviewAndContent(contentRaw)

	previewHtml, err := renderMarkdown(preview)
	if err != nil {
		return nil, err
	}

	contentHtml, err := renderMarkdown(content)
	if err != nil {
		return nil, err
	}

	return &Document{Meta: meta, Preview: string(previewHtml), Content: string(contentHtml)}, nil
}

func splitMetaAndContent(input []byte) ([]byte, []byte) {
	retval := strings.SplitN(string(input), "\n\n", 2)

	meta := []byte(retval[0])
	content := []byte(retval[1])

	return meta, content
}

func splitPreviewAndContent(input []byte) ([]byte, []byte) {
	retval := strings.SplitN(string(input), "\n[cut]\n", 2)

	var preview []byte
	var content []byte

	if len(retval) < 2 {
		preview = nil
		content = []byte(retval[0])
	} else {
		preview = []byte(retval[0])
		content = []byte(retval[1])
	}

	return preview, content
}

func parseMeta(input []byte) map[string]string {
	meta := make(map[string]string)

	pairs := strings.Split(string(input), "\n")

	for _, pair := range pairs {
		values := strings.SplitN(pair, ":", 2)
		meta[strings.TrimSpace(values[0])] = strings.TrimSpace(values[1])
	}

	return meta
}

func renderMarkdown(input []byte) ([]byte, error) {
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS

	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_XHTML

	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")
	var output []byte = blackfriday.Markdown(input, renderer, extensions)

	return output, nil
}
