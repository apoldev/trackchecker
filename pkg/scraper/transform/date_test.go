package transform_test

import (
	"github.com/apoldev/trackchecker/pkg/scraper/transform"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransformDate(t *testing.T) {

	cases := []struct {
		name   string
		data   string
		expect string
	}{
		{
			name:   "format 1",
			data:   "2024-02-05     10:00:00",
			expect: "2024-02-05T10:00:00Z",
		},
		{
			name:   "date format 2",
			data:   "January 4, 2024 1:25 pm",
			expect: "2024-01-04T13:25:00Z",
		},
		{
			name:   "unix data",
			data:   "/Date(1706534580000+0100)/",
			expect: "2024-01-29T13:23:00Z",
		},
		{

			data:   "/Date(1706534580+0100)/",
			expect: "2024-01-29T13:23:00Z",
		},
		{
			data:   "1706534520",
			expect: "2024-01-29T13:22:00Z",
		},

		{
			data:   "1706534+0100",
			expect: "1706534+0100",
		},
	}

	for _, c := range cases {
		d := transform.Date(c.data)
		require.Equal(t, c.expect, d)

	}
}
