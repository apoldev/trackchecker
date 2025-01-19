package document_test

import (
	"github.com/apoldev/trackchecker/pkg/scraper/document"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJSONXpathDoc_FindOne(t *testing.T) {

	cases := []struct {
		path        string
		want        interface{}
		expectError error
	}{
		{
			path: "//baz",
			want: "qux",
		},

		{
			path: "//foo/array",
			want: []interface{}{"apple", "orange", "pear"},
		},

		{
			path:        "",
			expectError: document.ErrInvalidQuery,
		},

		{
			path:        "[[[",
			expectError: document.ErrInvalidQuery,
		},

		{
			path: "//foo/array/*[1]",
			want: "apple",
		},

		{
			path: "concat(//foo/array/*[1], //foo/array/*[2])",
			want: "appleorange",
		},
	}

	doc, err := document.NewJSONXpath(JSONDATA)
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

func TestJSONXpathDoc_FindAll(t *testing.T) {

	cases := []struct {
		path string
		want []interface{}
	}{
		{
			path: "//baz",
			want: []interface{}{"qux"},
		},

		{
			path: "//bar",
			want: []interface{}{
				map[string]interface{}{
					"baz": "qux",
				},
			},
		},

		{
			path: "//data/gggg",
			want: []interface{}{},
		},

		{
			path: "//code/..",
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
			path: "//data/*",
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
			path: "[broken",
			want: []interface{}{},
		},
	}

	doc, err := document.NewJSONXpath(JSONDATA)
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
