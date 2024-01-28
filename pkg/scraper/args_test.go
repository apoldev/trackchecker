package scraper

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestResultBuilder(t *testing.T) {
	b := NewResultBuilder()

	assert.Equal(t, b.GetData(), []byte("{}"))
	assert.Equal(t, b.GetString(), "{}")

	b.Set("a.0.b", 777)

	assert.Equal(t, b.GetString(), `{"a":[{"b":777}]}`)

	// error case
	b.Set("a.0.b", json.RawMessage{0, 1, 2})
	assert.Equal(t, b.GetString(), `{"a":[{"b":777}]}`)
}
