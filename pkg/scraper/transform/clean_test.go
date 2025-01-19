package transform

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClean(t *testing.T) {

	cases := []struct {
		name   string
		data   string
		expect string
	}{
		{
			data:   ",  ",
			expect: "",
		},
		{
			data:   "/New York./",
			expect: "New York",
		},

		{
			data:   "Delivered",
			expect: "Delivered",
		},
	}

	for _, c := range cases {
		d := Clean(c.data)
		require.Equal(t, c.expect, d)
	}
}
