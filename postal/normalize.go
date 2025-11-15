package postal

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// NormalizeOptions configures string normalization behavior
type NormalizeOptions struct {
	Lowercase           bool
	RemoveAccents       bool
	ExpandNumex         bool
	LatinASCII          bool
	Transliterate       bool
	StripPunctuation    bool
	ReplaceNumericHyphens bool
}

// DefaultNormalizeOptions returns default normalization options
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		Lowercase:           true,
		RemoveAccents:       true,
		ExpandNumex:         false,
		LatinASCII:          true,
		Transliterate:       true,
		StripPunctuation:    false,
		ReplaceNumericHyphens: true,
	}
}

// Normalize normalizes a string according to the provided options
func Normalize(input string, options NormalizeOptions) string {
	if input == "" {
		return ""
	}

	result := input

	// Lowercase
	if options.Lowercase {
		result = strings.ToLower(result)
	}

	// Remove accents (NFD decomposition + remove combining marks)
	if options.RemoveAccents {
		result = removeAccents(result)
	}

	// Strip punctuation
	if options.StripPunctuation {
		result = stripPunctuation(result)
	}

	// Replace numeric hyphens (e.g., "123-45" stays, but "foo-bar" becomes "foo bar")
	if options.ReplaceNumericHyphens {
		result = replaceNumericHyphens(result)
	}

	// Normalize whitespace
	result = normalizeWhitespace(result)

	return strings.TrimSpace(result)
}

// removeAccents removes accents from characters
func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

// stripPunctuation removes punctuation characters
func stripPunctuation(s string) string {
	var builder strings.Builder
	for _, r := range s {
		if !unicode.IsPunct(r) {
			builder.WriteRune(r)
		} else {
			builder.WriteRune(' ')
		}
	}
	return builder.String()
}

// replaceNumericHyphens replaces hyphens between non-digits with spaces
func replaceNumericHyphens(s string) string {
	runes := []rune(s)
	result := make([]rune, 0, len(runes))

	for i := 0; i < len(runes); i++ {
		if runes[i] == '-' {
			// Check if hyphen is between digits
			prevIsDigit := i > 0 && unicode.IsDigit(runes[i-1])
			nextIsDigit := i < len(runes)-1 && unicode.IsDigit(runes[i+1])

			if prevIsDigit && nextIsDigit {
				result = append(result, '-')
			} else {
				result = append(result, ' ')
			}
		} else {
			result = append(result, runes[i])
		}
	}

	return string(result)
}

// normalizeWhitespace replaces multiple spaces with single space
func normalizeWhitespace(s string) string {
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

// NormalizeString is a convenience function using default options
func NormalizeString(input string) string {
	return Normalize(input, DefaultNormalizeOptions())
}
