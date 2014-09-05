package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderMarkdown(t *testing.T) {
	output, _ := renderMarkdown([]byte("Hello"))

	assert.Equal(t, "<p>Hello</p>\n", string(output))
}

func TestSplitMetaAndContent(t *testing.T) {
	meta, content := splitMetaAndContent([]byte("Title: Foo\n\nHello there!"))

	assert.Equal(t, "Title: Foo", string(meta))
	assert.Equal(t, "Hello there!", string(content))
}

func TestSplitMetaAndContentWithoutMeta(t *testing.T) {
	meta, content := splitMetaAndContent([]byte("Hello there!"))

	assert.Equal(t, "", string(meta))
	assert.Equal(t, "Hello there!", string(content))
}

func TestSplitPreviewAndContent(t *testing.T) {
	preview, content := splitPreviewAndContent([]byte("foo\n[cut]\nbar"))

	assert.Equal(t, "foo", string(preview))
	assert.Equal(t, "bar", string(content))
}

func TestSplitPreviewAndContentWithoutPreview(t *testing.T) {
	preview, content := splitPreviewAndContent([]byte("foo\nbar"))

	assert.Equal(t, "", string(preview))
	assert.Equal(t, "foo\nbar", string(content))
}

func TestParseMeta(t *testing.T) {
	meta := parseMeta([]byte("Foo: bar\nBar : baz"))

	assert.Equal(t, "bar", meta["Foo"])
	assert.Equal(t, "baz", meta["Bar"])
}

func TestParseEmptyMeta(t *testing.T) {
	meta := parseMeta([]byte(""))

	assert.Equal(t, 0, len(meta))
}

func TestParseInvalidMeta(t *testing.T) {
	meta := parseMeta([]byte("Foo"))

	assert.Equal(t, 0, len(meta))
}

func TestParseMetaWithEmptyValue(t *testing.T) {
	meta := parseMeta([]byte("Foo:"))

	assert.Equal(t, "", meta["Foo"])
}

func TestParseFileName(t *testing.T) {
	date, fileName := parseFileName("20110112-filename.markdown")

	assert.Equal(t, 2011, date["year"])
	assert.Equal(t, 1, date["month"])
	assert.Equal(t, 12, date["day"])
	assert.Equal(t, "filename.markdown", fileName)
}

func TestParseFileNameWithoutDate(t *testing.T) {
	date, fileName := parseFileName("filename.markdown")

	assert.Nil(t, date)
	assert.Equal(t, "filename.markdown", fileName)
}

func TestParseFileNameWithInvalidDate(t *testing.T) {
	date, fileName := parseFileName("99999999-filename.markdown")

	assert.Nil(t, date)
	assert.Equal(t, "99999999-filename.markdown", fileName)
}

func TestParseInput(t *testing.T) {
	document, _ := parseContent([]byte("Foo: bar\nBar : baz\n\npreview\n[cut]\ncontent"))

	assert.Equal(t, "bar", document.Meta["Foo"])
	assert.Equal(t, "baz", document.Meta["Bar"])
	assert.Equal(t, "<p>preview</p>\n", string(document.Preview))
	assert.Equal(t, "<p>content</p>\n", string(document.Content))
}

func TestParseDate(t *testing.T) {
	date, _ := parseDate("20111213")

	assert.Equal(t, 2011, date["year"])
	assert.Equal(t, 12, date["month"])
	assert.Equal(t, 13, date["day"])
}

func TestParseInvalidDate(t *testing.T) {
	date, err := parseDate("99999999-hello.there")

	assert.Nil(t, date)
	assert.NotNil(t, err)
}

func TestParseDocument(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	filename := "20110113-article.markdown"

	input := "Title: Foo\n\nHi there\n[cut]\nAnd here"
	ioutil.WriteFile(tempDir+filename, []byte(input), 0644)

	dm := DocumentManager{Root: tempDir}
	document, _ := dm.parseDocument(filename)

	assert.Equal(t, tempDir+filename, document.Path)
	assert.Equal(t, "article", document.Slug)
	assert.Equal(t, "2011", document.Created["Year"])
	assert.Equal(t, "01", document.Created["Month"])
	assert.Equal(t, "Foo", document.Meta["Title"])
	assert.Equal(t, "<p>Hi there</p>\n", document.Preview)
	assert.Equal(t, "<p>And here</p>\n", document.Content)

	os.RemoveAll(tempDir)
}

func TestParseDocumentWithUnknownFile(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	dm := DocumentManager{Root: tempDir}
	document, _ := dm.parseDocument("unknown-file")

	assert.Nil(t, document)

	os.RemoveAll(tempDir)
}

func TestParseDocumentWithUnreadableFile(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	filename := "file"
	ioutil.WriteFile(tempDir+filename, []byte{}, 0000)

	dm := DocumentManager{Root: tempDir}
	document, _ := dm.parseDocument(filename)

	assert.Nil(t, document)

	os.RemoveAll(tempDir)
}

func TestParseDocuments(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	filename := "20110113-article.markdown"

	input := "Title: Foo\n\nHi there\n[cut]\nAnd here"
	ioutil.WriteFile(tempDir+filename, []byte(input), 0644)

	ioutil.WriteFile(tempDir+".ignore-me", []byte{}, 0644)
	ioutil.WriteFile(tempDir+"ignore-me.foo", []byte{}, 0644)

	dm := DocumentManager{Root: tempDir}
	documents, _ := dm.parseDocuments()

	assert.Equal(t, 1, len(documents))
	assert.Equal(t, "Foo", documents[0].Meta["Title"])

	os.RemoveAll(tempDir)
}

func TestLoadDocumentBySlug(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	filename := "20110113-article.markdown"
	input := "Title: Foo\n\nHi there\n[cut]\nAnd here"
	ioutil.WriteFile(tempDir+filename, []byte(input), 0644)

	filename = "20110113-other-article.markdown"
	input = "Title: Bar\n\nHi there\n[cut]\nAnd here"
	ioutil.WriteFile(tempDir+filename, []byte(input), 0644)

	dm := DocumentManager{Root: tempDir}
	document, _ := dm.LoadDocumentBySlug("article")

	assert.Equal(t, "Foo", document.Meta["Title"])

	os.RemoveAll(tempDir)
}

func TestLoadDocumentBySlugAndDate(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	filename := "article.markdown"
	ioutil.WriteFile(tempDir+filename, []byte{}, 0644)

	filename = "20110113-another.markdown"
	ioutil.WriteFile(tempDir+filename, []byte{}, 0644)

	filename = "20110113-article.markdown"
	input := "Title: Foo\n\nHi there\n[cut]\nAnd here"
	ioutil.WriteFile(tempDir+filename, []byte(input), 0644)

	filename = "20110213-article.markdown"
	input = "Title: Bar\n\nHi there\n[cut]\nAnd here"
	ioutil.WriteFile(tempDir+filename, []byte(input), 0644)

	dm := DocumentManager{Root: tempDir}
	document, _ := dm.LoadDocumentBySlugAndDate("article", "2011", "02")

	assert.NotNil(t, document)
	assert.Equal(t, "Bar", document.Meta["Title"])

	os.RemoveAll(tempDir)
}

func TestLoadDocumentBySlugAndDateNothingFoundBySlug(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	filename := "20110113-article.markdown"
	input := "Title: Foo\n\nHi there\n[cut]\nAnd here"
	ioutil.WriteFile(tempDir+filename, []byte(input), 0644)

	dm := DocumentManager{Root: tempDir}
	document, _ := dm.LoadDocumentBySlugAndDate("something-else", "2011", "02")

	assert.Nil(t, document)

	os.RemoveAll(tempDir)
}

func TestLoadDocumentBySlugAndDateNothingFoundByDate(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	filename := "20110113-article.markdown"
	input := "Title: Foo\n\nHi there\n[cut]\nAnd here"
	ioutil.WriteFile(tempDir+filename, []byte(input), 0644)

	dm := DocumentManager{Root: tempDir}
	document, _ := dm.LoadDocumentBySlugAndDate("article", "2011", "02")

	assert.Nil(t, document)

	os.RemoveAll(tempDir)
}

func TestLoadDocuments(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	for _, filename := range []string{
		"20110113-very-old.markdown",
		"20110114-old.markdown",
		"20110115-new.markdown",
	} {
		ioutil.WriteFile(tempDir+filename, []byte{}, 0644)
	}

	dm := DocumentManager{Root: tempDir}
	documents, _ := dm.LoadDocuments(0, "")

	assert.Equal(t, 3, len(documents))
	assert.Equal(t, "new", documents[0].Slug)
	assert.Equal(t, "old", documents[1].Slug)
	assert.Equal(t, "very-old", documents[2].Slug)

	os.RemoveAll(tempDir)
}

func TestLoadDocumentsWithLimit(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	for _, filename := range []string{
		"20110113-very-old.markdown",
		"20110114-old.markdown",
		"20110115-new.markdown",
	} {
		ioutil.WriteFile(tempDir+filename, []byte{}, 0644)
	}

	dm := DocumentManager{Root: tempDir}
	documents, _ := dm.LoadDocuments(2, "")

	assert.Equal(t, 2, len(documents))
	assert.Equal(t, "new", documents[0].Slug)
	assert.Equal(t, "old", documents[1].Slug)

	os.RemoveAll(tempDir)
}

func TestLoadDocumentsWithLimitTooBig(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	for _, filename := range []string{
		"20110113-very-old.markdown",
		"20110114-old.markdown",
		"20110115-new.markdown",
	} {
		ioutil.WriteFile(tempDir+filename, []byte{}, 0644)
	}

	dm := DocumentManager{Root: tempDir}
	documents, _ := dm.LoadDocuments(100, "")

	assert.Equal(t, 3, len(documents))

	os.RemoveAll(tempDir)
}

func TestLoadDocumentsWithLimitAndOffset(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	for _, filename := range []string{
		"20110113-very-old.markdown",
		"20110114-old.markdown",
		"20110115-new.markdown",
	} {
		ioutil.WriteFile(tempDir+filename, []byte{}, 0644)
	}

	dm := DocumentManager{Root: tempDir}
	documents, _ := dm.LoadDocuments(2, "20110114")

	assert.Equal(t, 2, len(documents))
	assert.Equal(t, "old", documents[0].Slug)
	assert.Equal(t, "very-old", documents[1].Slug)

	os.RemoveAll(tempDir)
}

func TestNextPageOffset(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	for _, filename := range []string{
		"20110113-very-old.markdown",
		"20110114-old.markdown",
		"20110115-new.markdown",
	} {
		ioutil.WriteFile(tempDir+filename, []byte{}, 0644)
	}

	dm := DocumentManager{Root: tempDir}

	nextOffset := dm.NextPageOffset(2, "")

	assert.Equal(t, "20110113", nextOffset)

	os.RemoveAll(tempDir)
}

func TestNextPageOffsetNotEnough(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	for _, filename := range []string{
		"20110113-very-old.markdown",
		"20110114-old.markdown",
		"20110115-new.markdown",
	} {
		ioutil.WriteFile(tempDir+filename, []byte{}, 0644)
	}

	dm := DocumentManager{Root: tempDir}

	nextOffset := dm.NextPageOffset(5, "")

	assert.Equal(t, "", nextOffset)

	os.RemoveAll(tempDir)
}

func TestNextPageOffsetWithOffset(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	for _, filename := range []string{
		"20110112-very-very-old.markdown",
		"20110113-very-old.markdown",
		"20110114-old.markdown",
		"20110115-new.markdown",
	} {
		ioutil.WriteFile(tempDir+filename, []byte{}, 0644)
	}

	dm := DocumentManager{Root: tempDir}

	nextOffset := dm.NextPageOffset(2, "20110114")

	assert.Equal(t, "20110112", nextOffset)

	os.RemoveAll(tempDir)
}

func TestPrevPageOffsetNoOffset(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	for _, filename := range []string{
		"20110112-very-very-old.markdown",
		"20110113-very-old.markdown",
		"20110114-old.markdown",
		"20110115-new.markdown",
	} {
		ioutil.WriteFile(tempDir+filename, []byte{}, 0644)
	}

	dm := DocumentManager{Root: tempDir}

	prevOffset := dm.PrevPageOffset(2, "")

	assert.Equal(t, "", prevOffset)

	os.RemoveAll(tempDir)
}

func TestPrevPageOffsetNotEnough(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	for _, filename := range []string{
		"20110112-very-very-old.markdown",
		"20110113-very-old.markdown",
		"20110114-old.markdown",
		"20110115-new.markdown",
	} {
		ioutil.WriteFile(tempDir+filename, []byte{}, 0644)
	}

	dm := DocumentManager{Root: tempDir}

	prevOffset := dm.PrevPageOffset(2, "20110114")

	assert.Equal(t, "20110115", prevOffset)

	os.RemoveAll(tempDir)
}

func TestPrevPageOffsetNoWhere(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	for _, filename := range []string{
		"20110112-very-very-old.markdown",
		"20110113-very-old.markdown",
		"20110114-old.markdown",
		"20110115-new.markdown",
	} {
		ioutil.WriteFile(tempDir+filename, []byte{}, 0644)
	}

	dm := DocumentManager{Root: tempDir}

	prevOffset := dm.PrevPageOffset(2, "20110115")

	assert.Equal(t, "", prevOffset)

	os.RemoveAll(tempDir)
}

func TestPrevPageOffset(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	for _, filename := range []string{
		"20110112-very-very-old.markdown",
		"20110113-very-old.markdown",
		"20110114-old.markdown",
		"20110115-new.markdown",
	} {
		ioutil.WriteFile(tempDir+filename, []byte{}, 0644)
	}

	dm := DocumentManager{Root: tempDir}

	prevOffset := dm.PrevPageOffset(2, "20110112")

	assert.Equal(t, "20110114", prevOffset)

	os.RemoveAll(tempDir)
}

func TestNewerDocumentNotFound(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	for _, filename := range []string{
		"20110112-very-very-old.markdown",
		"20110113-very-old.markdown",
		"20110114-old.markdown",
		"20110115-new.markdown",
	} {
		ioutil.WriteFile(tempDir+filename, []byte{}, 0644)
	}

	dm := DocumentManager{Root: tempDir}

	got := dm.NewerDocument(&Document{Path: "foobar"})

	assert.Nil(t, got)

	os.RemoveAll(tempDir)
}

func TestNewerDocumentNo(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	for _, filename := range []string{
		"20110112-very-very-old.markdown",
		"20110113-very-old.markdown",
		"20110114-old.markdown",
		"20110115-new.markdown",
	} {
		ioutil.WriteFile(tempDir+filename, []byte{}, 0644)
	}

	dm := DocumentManager{Root: tempDir}

	got := dm.NewerDocument(&Document{Path: tempDir + "20110115-new.markdown"})

	assert.Nil(t, got)

	os.RemoveAll(tempDir)
}

func TestNewerDocument(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	for _, filename := range []string{
		"20110112-very-very-old.markdown",
		"20110113-very-old.markdown",
		"20110114-old.markdown",
		"20110115-new.markdown",
	} {
		ioutil.WriteFile(tempDir+filename, []byte{}, 0644)
	}

	dm := DocumentManager{Root: tempDir}

	got := dm.NewerDocument(&Document{Path: tempDir + "20110114-old.markdown"})

	assert.Equal(t, "new", got.Slug)

	os.RemoveAll(tempDir)
}

func TestOlderDocumentNotFound(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	for _, filename := range []string{
		"20110112-very-very-old.markdown",
		"20110113-very-old.markdown",
		"20110114-old.markdown",
		"20110115-new.markdown",
	} {
		ioutil.WriteFile(tempDir+filename, []byte{}, 0644)
	}

	dm := DocumentManager{Root: tempDir}

	got := dm.OlderDocument(&Document{Path: "foobar"})

	assert.Nil(t, got)

	os.RemoveAll(tempDir)
}

func TestOlderDocumentNo(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	for _, filename := range []string{
		"20110112-very-very-old.markdown",
		"20110113-very-old.markdown",
		"20110114-old.markdown",
		"20110115-new.markdown",
	} {
		ioutil.WriteFile(tempDir+filename, []byte{}, 0644)
	}

	dm := DocumentManager{Root: tempDir}

	got := dm.OlderDocument(&Document{Path: tempDir + "20110112-very-very.markdown"})

	assert.Nil(t, got)

	os.RemoveAll(tempDir)
}

func TestOlderDocument(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	tempDir = tempDir + "/"

	for _, filename := range []string{
		"20110112-very-very-old.markdown",
		"20110113-very-old.markdown",
		"20110114-old.markdown",
		"20110115-new.markdown",
	} {
		ioutil.WriteFile(tempDir+filename, []byte{}, 0644)
	}

	dm := DocumentManager{Root: tempDir}

	got := dm.OlderDocument(&Document{Path: tempDir + "20110113-very-old.markdown"})

	assert.Equal(t, "very-very-old", got.Slug)

	os.RemoveAll(tempDir)
}

func TestDocumentCreated(t *testing.T) {
	document := Document{Created: map[string]string{"Year": "2011", "Month": "02", "Day": "13"}}

	fmt.Println(document.Created)
	assert.Equal(t, "20110213", document.Created.String())
}
