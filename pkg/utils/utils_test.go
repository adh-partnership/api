package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayContains(t *testing.T) {
	assert.True(t, ArrayContains([]string{"a", "b", "c"}, "b"))
	assert.False(t, ArrayContains([]string{"a", "b", "c"}, "d"))
}
