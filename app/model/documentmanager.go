package model

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/russross/blackfriday"
)

type Document struct {
	Meta    map[string]string
	Preview []byte
	Content []byte
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

	return &Document{Meta: meta, Preview: previewHtml, Content: contentHtml}, nil
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
		preview = []byte(retval[0])
		content = preview
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
