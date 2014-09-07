package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

type DateValue int
type Date map[string]DateValue

func (v DateValue) String() string {
	return strconv.Itoa(int(v))
}
func (c Date) String() string {
	return strconv.Itoa(int(c["Year"])) +
		fmt.Sprintf("%02d", int(c["Month"])) +
		fmt.Sprintf("%02d", int(c["Day"]))
}

type Document struct {
	Path        string
	Slug        string
	Tags        []string
	Created     Date
	Meta        map[string]string
	Preview     string
	PreviewLink string
	Content     string
}

type DocumentManager struct {
	Root string
}

func (dm *DocumentManager) LoadDocumentBySlug(slug string) (*Document, error) {
	documents, err := dm.parseDocuments()
	if err != nil {
		return nil, err
	}

	for _, document := range documents {
		if document.Slug == slug {
			return document, nil
		}
	}

	return nil, nil
}

func (dm *DocumentManager) LoadDocumentBySlugAndDate(slug string, year string, month string) (*Document, error) {
	documents, err := dm.parseDocuments()
	if err != nil {
		return nil, err
	}

	for _, document := range documents {
		if document.Slug != slug {
			continue
		}

		if year == document.Created["Year"].String() &&
			fmt.Sprintf("%02d", month) == fmt.Sprintf("%02d", document.Created["Month"].String()) {
			return document, nil
		}
	}

	return nil, nil
}

func (dm *DocumentManager) LoadDocumentsByTag(tag string) ([]*Document, error) {
	documents, err := dm.LoadDocuments(0, "")
	if err != nil {
		return nil, err
	}

	taggedDocuments := []*Document{}
	for _, document := range documents {
		for _, t := range document.Tags {
			if t == tag {
				taggedDocuments = append(taggedDocuments, document)
				break
			}
		}
	}

	return taggedDocuments, nil
}

func (dm *DocumentManager) LoadDocuments(limit int, offset string) ([]*Document, error) {
	documents, err := dm.parseDocuments()
	if err != nil {
		return nil, err
	}

	if len(offset) > 0 {
		var index int
		for i, document := range documents {
			currentTimestamp := document.Created.String()
			if currentTimestamp > offset {
				continue
			}

			index = i
			break
		}

		if index > 0 {
			documents = documents[index:]
		}
	}

	if limit > 0 && limit < len(documents) {
		documents = documents[0:limit]
	}

	return documents, nil
}

func (dm *DocumentManager) NextPageOffset(limit int, offset string) string {
	if limit <= 0 {
		return ""
	}

	documents, _ := dm.LoadDocuments(limit+1, offset)
	if len(documents) != limit+1 {
		return ""
	}

	last := documents[len(documents)-1]
	return last.Created.String()
}

func (dm *DocumentManager) PrevPageOffset(limit int, offset string) string {
	if limit <= 0 || len(offset) == 0 {
		return ""
	}

	documents, _ := dm.LoadDocuments(0, "")
	var index int
	for i, document := range documents {
		currentTimestamp := document.Created.String()

		if currentTimestamp <= offset {
			index = i
			break
		}
	}

	if index == 0 {
		return ""
	}

	index = index - limit
	if index < 0 {
		index = 0
	}

	offsetted := documents[index]
	return offsetted.Created.String()
}

func (dm *DocumentManager) NewerDocument(document *Document) *Document {
	documents, _ := dm.LoadDocuments(0, "")

	var index int
	for i, d := range documents {
		if document.Path == d.Path {
			index = i
			break
		}
	}

	if index == 0 {
		return nil
	}

	return documents[index-1]
}

func (dm *DocumentManager) OlderDocument(document *Document) *Document {
	documents, _ := dm.LoadDocuments(0, "")

	var index int
	for i, d := range documents {
		index = i
		if document.Path == d.Path {
			break
		}
	}

	if index+1 >= len(documents) {
		return nil
	}

	return documents[index+1]
}

func (dm *DocumentManager) parseDocument(name string) (*Document, error) {
	path := dm.Root + name

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	input, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	document, err := parseContent(input)

	if err != nil {
		return nil, err
	}

	date, fileName := parseFileName(name)
	document.Created = date

	ext := filepath.Ext(fileName)
	slug := strings.TrimSuffix(fileName, ext)

	document.Path = path
	document.Slug = slug

	return document, nil
}

func (dm *DocumentManager) parseDocuments() ([]*Document, error) {
	files, err := listFiles(dm.Root)
	if err != nil {
		return nil, err
	}

	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	var documents []*Document
	for _, name := range files {
		document, err := dm.parseDocument(name)
		if err != nil {
			continue
		}

		documents = append(documents, document)
	}

	return documents, nil
}

func listFiles(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var documents []string
	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		if !strings.HasSuffix(f.Name(), ".markdown") {
			continue
		}

		documents = append(documents, f.Name())
	}

	return documents, nil
}

func parseFileName(input string) (Date, string) {
	dateRe := regexp.MustCompile("^[0-9]{8}-")

	if dateRe.MatchString(input) {
		parts := strings.SplitN(input, "-", 2)

		if len(parts) == 2 {
			date, _ := parseDate(parts[0])

			if date != nil {
				return date, parts[1]
			}
		}
	}

	return nil, input
}

func parseDate(input string) (Date, error) {
	t, err := time.Parse("20060102", input)
	if err != nil {
		return nil, err
	}

	date := Date{"Year": DateValue(t.Year()),
		"Month": DateValue(t.Month()),
		"Day":   DateValue(t.Day())}
	return date, nil
}

func parseContent(input []byte) (*Document, error) {
	metaRaw, contentRaw := splitMetaAndContent(input)
	meta := parseMeta(metaRaw)
	preview, previewLink, content := splitPreviewAndContent(contentRaw)

	previewHtml, err := renderMarkdown(preview)
	if err != nil {
		return nil, err
	}

	contentHtml, err := renderMarkdown(content)
	if err != nil {
		return nil, err
	}

	tags := []string{}
	if len(meta["Tags"]) > 0 {
		re := regexp.MustCompile("\\s*,+\\s*")
		tags = re.Split(string(meta["Tags"]), -1)
	}

	return &Document{
		Meta:        meta,
		Tags:        tags,
		Preview:     string(previewHtml),
		PreviewLink: string(previewLink),
		Content:     string(contentHtml)}, nil
}

func splitMetaAndContent(input []byte) (meta, content []byte) {
	retval := strings.SplitN(string(input), "\n\n", 2)

	if len(retval) == 2 {
		meta = []byte(retval[0])
		content = []byte(retval[1])
	} else {
		meta = []byte{}
		content = input
	}

	return
}

func splitPreviewAndContent(input []byte) ([]byte, []byte, []byte) {
	re := regexp.MustCompile("\n\\[cut\\](?: +([^\r\n]+))?\r?\n")

	retval := re.Split(string(input), -1)

	var preview []byte
	var previewLink []byte
	var content []byte

	if len(retval) < 2 {
		preview = nil
		content = []byte(retval[0])
	} else {
		submatch := re.FindStringSubmatch(string(input))
		if len(submatch) > 1 {
			previewLink = []byte(submatch[1])
		}

		preview = []byte(retval[0])
		content = []byte(retval[1])
	}

	return preview, previewLink, content
}

func parseMeta(input []byte) map[string]string {
	meta := map[string]string{}

	pairs := strings.Split(string(input), "\n")

	for _, pair := range pairs {
		values := strings.SplitN(pair, ":", 2)
		if len(values) == 2 {
			name := strings.TrimSpace(values[0])
			meta[name] = strings.TrimSpace(values[1])
		}
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
