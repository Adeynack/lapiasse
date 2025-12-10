package model

import (
	"bytes"
	"encoding/json/jsontext"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestID_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		id       ID
		expected string
	}{
		{
			name:     "zero value",
			id:       ID(0),
			expected: `"0"`,
		},
		{
			name:     "small positive number",
			id:       ID(42),
			expected: `"42"`,
		},
		{
			name:     "large positive number",
			id:       ID(9223372036854775807),
			expected: `"9223372036854775807"`,
		},
		{
			name:     "max uint64",
			id:       ID(18446744073709551615),
			expected: `"18446744073709551615"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			enc := jsontext.NewEncoder(&buf)
			err := tt.id.MarshalJSONTo(enc)
			require.NoError(t, err)

			result := buf.Bytes()
			assert.JSONEq(t, tt.expected, string(result))
		})
	}
}

func TestID_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ID
		wantErr  bool
	}{
		{
			name:     "zero value",
			input:    `"0"`,
			expected: ID(0),
			wantErr:  false,
		},
		{
			name:     "small positive number",
			input:    `"42"`,
			expected: ID(42),
			wantErr:  false,
		},
		{
			name:     "large positive number",
			input:    `"9223372036854775807"`,
			expected: ID(9223372036854775807),
			wantErr:  false,
		},
		{
			name:     "max uint64",
			input:    `"18446744073709551615"`,
			expected: ID(18446744073709551615),
			wantErr:  false,
		},
		{
			name:     "invalid - number instead of string",
			input:    `42`,
			expected: ID(0),
			wantErr:  true,
		},
		{
			name:     "invalid - negative number",
			input:    `"-1"`,
			expected: ID(0),
			wantErr:  true,
		},
		{
			name:     "invalid - non-numeric string",
			input:    `"abc"`,
			expected: ID(0),
			wantErr:  true,
		},
		{
			name:     "invalid - empty string",
			input:    `""`,
			expected: ID(0),
			wantErr:  true,
		},
		{
			name:     "invalid - overflow uint64",
			input:    `"18446744073709551616"`,
			expected: ID(0),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := jsontext.NewDecoder(bytes.NewReader([]byte(tt.input)))
			var id ID
			err := id.UnmarshalJSONFrom(dec)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, id)
			}
		})
	}
}

func TestID_MarshalUnmarshal_RoundTrip(t *testing.T) {
	tests := []struct {
		name string
		id   ID
	}{
		{name: "zero", id: ID(0)},
		{name: "one", id: ID(1)},
		{name: "small", id: ID(42)},
		{name: "medium", id: ID(1234567)},
		{name: "large", id: ID(9223372036854775807)},
		{name: "max", id: ID(18446744073709551615)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			var buf bytes.Buffer
			enc := jsontext.NewEncoder(&buf)
			err := tt.id.MarshalJSONTo(enc)
			require.NoError(t, err)

			marshaled := buf.Bytes()

			// Unmarshal
			dec := jsontext.NewDecoder(bytes.NewReader(marshaled))
			var result ID
			err = result.UnmarshalJSONFrom(dec)
			require.NoError(t, err)

			// Verify round trip
			assert.Equal(t, tt.id, result)
		})
	}
}
