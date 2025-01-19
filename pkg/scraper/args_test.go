package scraper

import (
	"encoding/json"
	"github.com/apoldev/trackchecker/pkg/scraper/transform"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestResultBuilder(t *testing.T) {
	b := NewResultBuilder()

	assert.Equal(t, []byte("{}"), b.GetData())
	assert.Equal(t, "{}", b.GetString())

	b.Set("a.0.b", 777)

	assert.Equal(t, `{"a":[{"b":777}]}`, b.GetString())

	// error case
	b.Set("a.0.b", json.RawMessage{0, 1, 2})
	assert.Equal(t, `{"a":[{"b":777}]}`, b.GetString())

	b.Set("a.1.b", "New York...")
	require.Equal(t, `{"a":[{"b":777},{"b":"New York..."}]}`, b.GetString())

	b.Set("a.1.b", "New York...", transform.Transformer{
		Type: transform.TypeClean,
	})
	require.Equal(t, `{"a":[{"b":777},{"b":"New York"}]}`, b.GetString())

	b.Set("a.2.b", "January 29, 2024 8:03 pm", transform.Transformer{
		Type: transform.TypeDate,
	})
	require.Equal(t, `{"a":[{"b":777},{"b":"New York"},{"b":"2024-01-29T20:03:00Z"}]}`, b.GetString())

	b.Set("a.2.b", "x.x", transform.Transformer{
		Type: transform.TypeReplaceString,
		Params: map[string]string{
			"old": ".",
			"new": "---",
		},
	}, transform.Transformer{Type: transform.TypeDate})
	require.Equal(t, `{"a":[{"b":777},{"b":"New York"},{"b":"x---x"}]}`, b.GetString())

	// test replace transformer with regexp
	b.Set("a.2.b", "abcdefg00932", transform.Transformer{
		Type: transform.TypeReplaceRegexp,
		Params: map[string]string{
			"regexp": "^(.*?)([0-9]+)$",
			"new":    "$2$1",
		},
	})
	require.Equal(t, `{"a":[{"b":777},{"b":"New York"},{"b":"00932abcdefg"}]}`, b.GetString())

}
