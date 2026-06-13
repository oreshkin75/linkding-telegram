package utils

import (
	"testing"

	"github.com/c2fo/testify/assert"
)

func TestParseURLs(t *testing.T) {
	assert.Equal(t, []string{"https://example.com", "http://example.org/path"}, ParseURLs("urls: https://example.com and http://example.org/path"))
}
