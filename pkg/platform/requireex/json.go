package requireex

import (
	"encoding/json/v2"
	"testing"

	"github.com/nsf/jsondiff"
	"github.com/stretchr/testify/require"
)

// NoJsonDiffFromStruct checks that the JSON marshalled from the given struct
// matches the expected JSON string, failing the test if there is a mismatch.
func NoJsonDiffFromStruct(t *testing.T, v any, expectedJSON string) {
	data, err := json.Marshal(v)
	require.NoError(t, err)

	NoJsonDiff(t, expectedJSON, string(data))
}

// NoJsonDiff checks that the actual JSON string matches the expected JSON string,
// failing the test if there is a mismatch.
func NoJsonDiff(t *testing.T, actualJSON string, expectedJSON string) {
	jsonDiffOpts := jsondiff.DefaultConsoleOptions()
	res, diff := jsondiff.Compare([]byte(expectedJSON), []byte(actualJSON), &jsonDiffOpts)
	if res != jsondiff.FullMatch {
		require.Failf(t, "JSON mismatch: %s", diff)
	}
}
