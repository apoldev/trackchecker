package transform_test

import (
	"github.com/apoldev/trackchecker/pkg/scraper/transform"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReplaceStr(t *testing.T) {
	cases := []struct {
		name        string
		str         string
		replaceFunc func(string, string, string) string
		params      []string
		expected    string
	}{
		{
			str:         "Hello, world!",
			replaceFunc: transform.ReplaceStr,
			params:      []string{"world", "universe"},
			expected:    "Hello, universe!",
		},

		{
			str:         "01234abc",
			replaceFunc: transform.ReplaceRegexp,
			params:      []string{"^([0-9]+)(.*?)$", "$2$1"},
			expected:    "abc01234",
		},

		{
			name:        "broken regexp",
			str:         "___",
			replaceFunc: transform.ReplaceRegexp,
			params:      []string{"[[", "$1"},
			expected:    "___",
		},
	}

	for _, c := range cases {
		d := c.replaceFunc(c.str, c.params[0], c.params[1])
		require.Equal(t, c.expected, d)
	}

}
