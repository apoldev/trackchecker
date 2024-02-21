package document_test

import (
	"github.com/apoldev/trackchecker/pkg/scraper/document"
	"github.com/stretchr/testify/require"
	"testing"
)

var XMLDATA = []byte(`<?xml version="1.0" encoding="utf-8"?>
<root>
<age>31</age>
<female>false</female>
<name>John</name>
<data>
<book>
<author>A</author>
<name>A.A</name>
</book>
<book>
<author>B</author>
<name>B.B</name>
</book>
</data>

<value>result<number>1000</number></value>

<a>
<x>1</x>
<x>2</x>
<x>3</x>
</a>

</root>`)

func TestXMLXpathDoc_FindOne(t *testing.T) {

	cases := []struct {
		path        string
		want        interface{}
		expectError error
	}{
		{
			path: "//age",
			want: "31",
		},

		{
			path:        "[brokenage",
			expectError: document.ErrInvalidQuery,
		},

		{
			path: "concat('1', //name[contains(text(), 'B')])",
			want: "1B.B",
		},

		{
			path: `//book/name[text()="A.A"]/../author`,
			want: "A",
		},

		{
			path: `//value/text()`,
			want: "result",
		},

		{
			path: `number(//value/number) * 2`,
			want: float64(2000),
		},
	}

	doc, err := document.NewXMLXpath(XMLDATA)
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

func TestXMLXpathDoc_FindAll(t *testing.T) {

	cases := []struct {
		path string
		want []interface{}
	}{
		{
			path: "//aa",
			want: []interface{}{},
		},
		{
			path: "[broken",
			want: []interface{}{},
		},

		{
			path: "//a/x",
			want: []interface{}{
				"1", "2", "3",
			},
		},

		{
			path: "//x[position()>1]",
			want: []interface{}{
				"2", "3",
			},
		},
	}

	doc, err := document.NewXMLXpath(XMLDATA)
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
