package models

import (
	"testing"
)

// test StringifiedInt unmarshalling
func TestUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name        string
		input       []byte
		expected    StringifiedInt
		expectError bool
	}{
		{
			name:        "Valid Stringified Number",
			input:       []byte(`"1234"`),
			expected:    1234,
			expectError: false,
		},
		{
			name:        "Valid Stringified Number",
			input:       []byte(`"1844674407370955161"`),
			expected:    1844674407370955161,
			expectError: false,
		},
		{
			name:        "Invalid Stringified Number",
			input:       []byte(`"abc"`),
			expectError: true,
		},
		{
			name:        "Invalid Stringified Number",
			input:       []byte(``),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var si StringifiedInt
			err := si.UnmarshalJSON(tc.input)

			if (err != nil) != tc.expectError {
				t.Errorf("Unexpected error status: %v, expectError: %v", err, tc.expectError)
			}

			if !tc.expectError && si != tc.expected {
				t.Errorf("Got %v, want %v", si, tc.expected)
			}
		})
	}
}
