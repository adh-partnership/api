package storage

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestGenerateURL(t *testing.T) {
	assert.Equal(t, "https://cdn.denartcc.org/uploads/test", GenerateURL("test"))
}

func TestGetSlugFromURL(t *testing.T) {
	assert.Equal(t, "test", GetSlugFromURL("https://cdn.denartcc.org/uploads/test"))
}
