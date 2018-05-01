package common

import (
	"testing"
)

func TestInt64Or(t *testing.T) {
	testCases := []struct {
		str      string
		def      int64
		expected int64
	}{
		{
			str:      "",
			def:      0,
			expected: 0,
		},
		{
			str:      "",
			def:      3,
			expected: 3,
		},
		{
			str:      "",
			def:      -99,
			expected: -99,
		},
		{
			str:      "x",
			def:      0,
			expected: 0,
		},
		{
			str:      "y",
			def:      4,
			expected: 4,
		},
		{
			str:      "z",
			def:      -100,
			expected: -100,
		},
		{
			str:      "4",
			def:      4,
			expected: 4,
		},
		{
			str:      "4",
			def:      -4,
			expected: 4,
		},
		{
			str:      "-4",
			def:      4,
			expected: -4,
		},
		{
			str:      "9000a",
			def:      -101,
			expected: -101,
		},
		{
			str:      "-40000000000",
			def:      0,
			expected: -40000000000,
		},
	}

	for i, testCase := range testCases {
		if actual := Int64Or(testCase.str, testCase.def); actual != testCase.expected {
			t.Errorf("[i=%v] Expected Int64Or(%q, %v)=%v but actual=%v", i, testCase.str, testCase.def, actual, testCase.expected)
		}
	}
}
