package document_test

import (
	"github.com/apoldev/trackchecker/pkg/scraper/document"
	"github.com/stretchr/testify/require"
	"testing"
)

var JSONDATA = []byte(`{
	"a": 1,
	"b": "2",
	"foo": {
		"bar": {
			"baz": "qux"
		},
		"array": ["apple", "orange", "pear"]
	},
	"data": [
	{
		"code": 1
	},
	{
		"code": 2
	}
	]
}`)

func TestJSONDoc_FindAll(t *testing.T) {
	t.Parallel()

	cases := []struct {
		path string
		want []interface{}
	}{

		{
			path: "foo.array",
			want: []interface{}{"apple", "orange", "pear"},
		},

		{
			path: "data",
			want: []interface{}{
				map[string]interface{}{
					"code": float64(1),
				},
				map[string]interface{}{
					"code": float64(2),
				},
			},
		},

		{
			path: "foo.bar",
			want: []interface{}{
				map[string]interface{}{
					"baz": "qux",
				},
			},
		},

		{
			path: "a",
			want: []interface{}{float64(1)},
		},

		{
			path: "b",
			want: []interface{}{"2"},
		},

		{
			path: "xxx",
			want: []interface{}{},
		},

		{
			path: "b.x.b.v",
			want: []interface{}{},
		},
	}

	doc, err := document.NewJSON(JSONDATA)
	require.NoError(t, err)

	for _, c := range cases {
		t.Run(c.path, func(t *testing.T) {
			docs := doc.FindAll(c.path)

			require.Len(t, docs, len(c.want))

			for i := range docs {
				require.Equal(t, c.want[i], docs[i].Value())
			}

		})
	}

}

func TestJSONDoc_FindOne(t *testing.T) {

	cases := []struct {
		path        string
		want        interface{}
		expectError error
	}{
		{
			path: "foo.bar.baz",
			want: "qux",
		},

		{
			path: "foo.array",
			want: []interface{}{"apple", "orange", "pear"},
		},

		{
			path:        "foo.unknown",
			expectError: document.ErrNotExists,
		},

		{
			path:        "broken query   ///",
			expectError: document.ErrNotExists,
		},
	}

	doc, err := document.NewJSON(JSONDATA)
	require.NoError(t, err)

	for _, c := range cases {
		t.Run(c.path, func(t *testing.T) {
			got, err := doc.FindOne(c.path)

			if c.expectError != nil {
				require.EqualError(t, err, c.expectError.Error())
				require.Empty(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.want, got.Value())
			}

		})
	}

}
