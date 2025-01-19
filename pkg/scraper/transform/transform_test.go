package transform_test

import (
	"encoding/json"
	"errors"
	"github.com/apoldev/trackchecker/pkg/scraper/transform"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestStruct struct {
	Data         int                    `json:"data"`
	Transformers transform.Transformers `json:"transformers"`
}

func TestTransformers_UnmarshalJSON(t *testing.T) {

	cases := []struct {
		data        string
		expected    transform.Transformers
		expectedErr error
	}{
		{
			data:     `{"data": 1, "transformers": [{"type": "clean"}, "date", {"type": "clean", "params": {"a": "b"}}]}`,
			expected: transform.Transformers{{Type: "clean"}, {Type: "date"}, {Type: "clean", Params: map[string]string{"a": "b"}}},
		},

		{
			data:     `{"data": 999, "transformers": [1, "clean", {"type":"date", "params": null}, {"type": "clean", "params": {"a": "b"}}]}`,
			expected: transform.Transformers{{Type: "clean"}, {Type: "date"}, {Type: "clean", Params: map[string]string{"a": "b"}}},
		},

		{
			data:        `{"data": 999, "transformers": 0}`,
			expectedErr: errors.New("cannot unmarshal number into Go"),
		},
	}

	for _, c := range cases {
		var tmp TestStruct
		err := json.Unmarshal([]byte(c.data), &tmp)
		if c.expectedErr != nil {
			require.ErrorContains(t, err, c.expectedErr.Error())
		} else {
			require.NoError(t, err)
			require.Equal(t, c.expected, tmp.Transformers)
		}

	}

}
