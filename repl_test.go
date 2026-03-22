package main

import "testing"

func TestCleanInput(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected []string
	}{
		"basic":              {input: "hello world", expected: []string{"hello", "world"}},
		"some caps":          {input: "HeLlo WoRld", expected: []string{"hello", "world"}},
		"all caps":           {input: "HELLO WORLD", expected: []string{"hello", "world"}},
		"leading space":      {input: " hello world", expected: []string{"hello", "world"}},
		"some nums and caps": {input: "HeLl0 W0r7d", expected: []string{"hell0", "w0r7d"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := cleanInput(tc.input)
			if len(actual) != len(tc.expected) {
				t.Fatalf("Mismatched lengths! Got: %v Expected: %v", len(actual), len(tc.expected))
			}
			for i := range actual {
				word := actual[i]
				expectedWord := tc.expected[i]
				if word != expectedWord {
					t.Fatalf("Mismatched words! Got: %v Expected: %v", word, expectedWord)
				}
			}
		})
	}
}
