package storage

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestGenerateURL(t *testing.T) {
	assert.Equal(t, "https://cdn.zdvartcc.org/test", GenerateURL("test"))
}

func TestGetSlugFromURL(t *testing.T) {
	assert.Equal(t, "test", GetSlugFromURL("https://cdn.zdvartcc.org/uploads/test"))
}
