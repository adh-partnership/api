package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidFormat(t *testing.T) {
	assert.True(t, IsValidFormat("text"))
	assert.True(t, IsValidFormat("json"))
	assert.False(t, IsValidFormat(""))
	assert.False(t, IsValidFormat("foo"))
}

func TestIsValidLogLevel(t *testing.T) {
	assert.True(t, IsValidLogLevel("trace"))
	assert.True(t, IsValidLogLevel("debug"))
	assert.True(t, IsValidLogLevel("info"))
	assert.True(t, IsValidLogLevel("warn"))
	assert.True(t, IsValidLogLevel("error"))
	assert.True(t, IsValidLogLevel("fatal"))
	assert.True(t, IsValidLogLevel("panic"))
	assert.False(t, IsValidLogLevel(""))
	assert.False(t, IsValidLogLevel("foo"))
}
