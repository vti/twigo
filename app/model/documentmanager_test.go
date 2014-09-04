package model

import (
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

func TestSplitPreviewAndContent(t *testing.T) {
	preview, content := splitPreviewAndContent([]byte("foo\n[cut]\nbar"))

	assert.Equal(t, "foo", string(preview))
	assert.Equal(t, "bar", string(content))
}

func TestSplitPreviewAndContentWithoutPreview(t *testing.T) {
	preview, content := splitPreviewAndContent([]byte("foo\nbar"))

	assert.Equal(t, "foo\nbar", string(preview))
	assert.Equal(t, "foo\nbar", string(content))
}

func TestParseMeta(t *testing.T) {
	meta := parseMeta([]byte("Foo: bar\nBar : baz"))

	assert.Equal(t, "bar", meta["Foo"])
	assert.Equal(t, "baz", meta["Bar"])
}

func TestParseInput(t *testing.T) {
	document, _ := parseInput([]byte("Foo: bar\nBar : baz\n\npreview\n[cut]\ncontent"))

	assert.Equal(t, "bar", document.Meta["Foo"])
	assert.Equal(t, "baz", document.Meta["Bar"])
	assert.Equal(t, "<p>preview</p>\n", string(document.Preview))
	assert.Equal(t, "<p>content</p>\n", string(document.Content))
}
