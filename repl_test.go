package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    " \t ",
			expected: []string{},
		},
		{
			input:    " GoLang Is AWESOME! ",
			expected: []string{"golang", "is", "awesome!"},
		},
		{
			input:    "BOOTDEV",
			expected: []string{"bootdev"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Error: Actual input (%v, Character Count: %v) does not match expected input (%v, Character Count: %v)", actual, len(actual), c.expected, len(c.expected))
			return
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Error: Mismatched word at index %v. %v is not %v. Original input: %v", i, word, expectedWord, c.input)
			}
		}
	}
}
