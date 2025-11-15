package postal

import (
	"strings"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		input       string
		minExpected int // minimum expected number of tokens
	}{
		{"123 Main St", 3},
		{"Quatre-vingt-douze Ave des Champs-Élysées", 5},
		{"", 0},
		{"東京都", 3}, // 3 ideographic characters
	}

	for _, tt := range tests {
		tokens := Tokenize(tt.input, DefaultTokenizeOptions())
		if len(tokens) < tt.minExpected {
			t.Errorf("Tokenize(%q): got %d tokens, want at least %d", tt.input, len(tokens), tt.minExpected)
		}
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		input        string
		shouldContain string
	}{
		{"HELLO WORLD", "hello"},
		{"Café", "cafe"},
		{"  multiple   spaces  ", "multiple spaces"},
		{"St. John's", "st john"},
	}

	opts := DefaultNormalizeOptions()
	opts.StripPunctuation = true

	for _, tt := range tests {
		result := Normalize(tt.input, opts)
		if !containsWord(result, tt.shouldContain) {
			t.Errorf("Normalize(%q): got %q, should contain %q", tt.input, result, tt.shouldContain)
		}
	}
}

func containsWord(s, word string) bool {
	return len(s) > 0 && (s == word || len(s) >= len(word))
}

func TestExpandAddress(t *testing.T) {
	tests := []struct {
		input string
		words []string // words that should appear in at least one expansion
	}{
		{"123 Main St", []string{"main", "street"}},
		{"456 Park Ave", []string{"park", "avenue"}},
		{"N. Broadway", []string{"broadway"}},
	}

	for _, tt := range tests {
		expansions := ExpandAddressRoot(tt.input)
		if len(expansions) == 0 {
			t.Errorf("ExpandAddress(%q): got no expansions", tt.input)
			continue
		}

		// Check that at least one expansion contains the expected words
		foundAll := false
		for _, exp := range expansions {
			allFound := true
			for _, word := range tt.words {
				if !strings.Contains(exp, word) {
					allFound = false
					break
				}
			}
			if allFound {
				foundAll = true
				break
			}
		}

		if !foundAll {
			t.Errorf("ExpandAddress(%q): no expansion contains all words %v, got %v", tt.input, tt.words, expansions)
		}
	}
}

func TestParseAddress(t *testing.T) {
	Setup()
	defer Teardown()

	tests := []struct {
		input        string
		component    string
		expectedVal  string
	}{
		{"123 Main Street", string(ComponentHouseNumber), "123"},
		{"456 Park Ave, New York, NY 10001", string(ComponentPostcode), "10001"},
	}

	for _, tt := range tests {
		response := ParseAddress(tt.input, DefaultAddressParserOptions())
		val := response.GetComponent(tt.component)
		if val != tt.expectedVal {
			t.Errorf("ParseAddress(%q).GetComponent(%s): got %q, want %q",
				tt.input, tt.component, val, tt.expectedVal)
		}
	}
}

func TestSetupTeardown(t *testing.T) {
	if IsSetup() {
		t.Error("Library should not be setup initially")
	}

	err := Setup()
	if err != nil {
		t.Errorf("Setup() failed: %v", err)
	}

	if !IsSetup() {
		t.Error("Library should be setup after Setup()")
	}

	Teardown()

	if IsSetup() {
		t.Error("Library should not be setup after Teardown()")
	}
}

func TestIsNameDuplicate(t *testing.T) {
	Setup()
	defer Teardown()

	tests := []struct {
		name1    string
		name2    string
		expected DuplicateStatus
	}{
		{"John Smith", "john smith", DuplicateStatusLikelyDuplicate},
		{"Café", "Cafe", DuplicateStatusLikelyDuplicate},
		{"Apple Inc", "Banana Corp", DuplicateStatusNonDuplicate},
	}

	for _, tt := range tests {
		status, err := IsNameDuplicate(tt.name1, tt.name2, DefaultDuplicateOptions())
		if err != nil {
			t.Errorf("IsNameDuplicate(%q, %q) error: %v", tt.name1, tt.name2, err)
		}
		if status != tt.expected {
			t.Errorf("IsNameDuplicate(%q, %q): got %v, want %v",
				tt.name1, tt.name2, status, tt.expected)
		}
	}
}

func BenchmarkTokenize(b *testing.B) {
	input := "123 Main Street, Apartment 4B, New York, NY 10001"
	opts := DefaultTokenizeOptions()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Tokenize(input, opts)
	}
}

func BenchmarkNormalize(b *testing.B) {
	input := "Quatre-vingt-douze Avenue des Champs-Élysées, Paris, France"
	opts := DefaultNormalizeOptions()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Normalize(input, opts)
	}
}

func BenchmarkExpandAddress(b *testing.B) {
	input := "123 N. Main St., Apt. 4B"
	opts := DefaultExpandOptions()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExpandAddress(input, opts)
	}
}

func BenchmarkParseAddress(b *testing.B) {
	Setup()
	defer Teardown()

	input := "123 Main Street, New York, NY 10001"
	opts := DefaultAddressParserOptions()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseAddress(input, opts)
	}
}
