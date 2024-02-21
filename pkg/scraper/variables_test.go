package scraper

import "testing"

func TestVariables_ReplaceStringFromVariables(t *testing.T) {
	cases := []struct {
		variable Variables
		value    string
		isErr    bool
		err      error
		wants    string
	}{
		{
			variable: Variables{
				"[track]": "000",
				"[yyy]":   "xxx",
			},
			value: "body=[track]",
			wants: "body=000",
		},
		{
			variable: Variables{
				"[track]": "",
			},
			isErr: true,
			value: "body=[track]",
		},
	}

	for i := range cases {
		c := cases[i]

		got := c.variable.ReplaceStringFromVariables(c.value)

		if got != c.wants && !c.isErr {
			t.Errorf("ReplaceStringFromVariables() == %q, want %q", got, c.wants)
		}

		if c.isErr && got == c.wants {
			t.Errorf("ReplaceStringFromVariables() == %q, want %q", got, c.wants)
		}
	}
}
