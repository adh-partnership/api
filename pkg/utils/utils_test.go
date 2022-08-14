package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayContains(t *testing.T) {
	assert.True(t, ArrayContains([]string{"a", "b", "c"}, "b"))
	assert.False(t, ArrayContains([]string{"a", "b", "c"}, "d"))
}

func TestStringToSlug(t *testing.T) {
	tests := []struct {
		Name           string
		Input          string
		ExpectedOutput string
	}{
		{
			Name:           "Empty",
			Input:          "",
			ExpectedOutput: "",
		},
		{
			Name:           "Short",
			Input:          "Hello",
			ExpectedOutput: "hello",
		},
		{
			Name:           "Long",
			Input:          "Hello, World!",
			ExpectedOutput: "hello-world",
		},
		{
			Name:           "Longer",
			Input:          "Hello, World! This is a very long string.",
			ExpectedOutput: "hello-world-this-is-a-very-long-string",
		},
		{
			Name: "Long with truncate",
			Input: "non odio euismod lacinia at quis risus sed vulputate odio ut enim blandit volutpat maecenas volutpat blandit aliquam etiam erat" +
				" velit scelerisque in dictum",
			ExpectedOutput: "non-odio-euismod-lacinia-at-quis-risus-sed-vulputate-odio-ut-enim-blandit-volutpat-maecenas-volutpa",
		},
		{
			Name:           "Test with special characters that should get filtered out",
			Input:          "Hello+!@%$(!% World",
			ExpectedOutput: "hello-world",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.ExpectedOutput, StringToSlug(test.Input))
		})
	}
}
