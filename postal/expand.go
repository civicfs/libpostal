package postal

import (
	"strings"
)

// ExpandOptions configures address expansion behavior
type ExpandOptions struct {
	Languages              []string
	AddressComponents      uint16
	LatinASCII             bool
	Transliterate          bool
	StripAccents           bool
	Lowercase              bool
	TrimString             bool
	ReplaceWordHyphens     bool
	DeleteWordHyphens      bool
	DeleteFinalPeriods     bool
	DeleteAcronymPeriods   bool
	DropEnglishPossessives bool
	DeleteApostrophes      bool
	ExpandNumex            bool
	RomanNumerals          bool
}

// DefaultExpandOptions returns default expansion options
func DefaultExpandOptions() ExpandOptions {
	return ExpandOptions{
		Languages:              nil, // Auto-detect
		AddressComponents:      0xFFFF,
		LatinASCII:             true,
		Transliterate:          true,
		StripAccents:           true,
		Lowercase:              true,
		TrimString:             true,
		ReplaceWordHyphens:     true,
		DeleteWordHyphens:      true,
		DeleteFinalPeriods:     true,
		DeleteAcronymPeriods:   true,
		DropEnglishPossessives: true,
		DeleteApostrophes:      true,
		ExpandNumex:            true,
		RomanNumerals:          true,
	}
}

// ExpandAddress expands an address into multiple normalized forms
// This is useful for search indexing and matching
func ExpandAddress(input string, options ExpandOptions) []string {
	if input == "" {
		return nil
	}

	// Start with the original input
	expansions := make(map[string]bool)
	expansions[input] = true

	// Normalize and create variations
	normalized := Normalize(input, NormalizeOptions{
		Lowercase:     options.Lowercase,
		RemoveAccents: options.StripAccents,
		LatinASCII:    options.LatinASCII,
		Transliterate: options.Transliterate,
	})
	expansions[normalized] = true

	// Expand abbreviations
	tokens := Tokenize(normalized, DefaultTokenizeOptions())
	expanded := expandAbbreviations(tokens)
	for _, exp := range expanded {
		expansions[exp] = true
	}

	// Handle apostrophes
	if options.DeleteApostrophes {
		withoutApostrophes := strings.ReplaceAll(normalized, "'", "")
		expansions[withoutApostrophes] = true
	}

	// Handle hyphens
	if options.DeleteWordHyphens {
		withoutHyphens := strings.ReplaceAll(normalized, "-", " ")
		withoutHyphens = normalizeWhitespace(withoutHyphens)
		expansions[withoutHyphens] = true
	}

	// Handle periods
	if options.DeleteFinalPeriods {
		withoutFinalPeriods := removeFinalPeriods(normalized)
		expansions[withoutFinalPeriods] = true
	}

	// Convert map to slice
	result := make([]string, 0, len(expansions))
	for exp := range expansions {
		if exp != "" {
			result = append(result, strings.TrimSpace(exp))
		}
	}

	return result
}

// expandAbbreviations expands common abbreviations in address tokens
func expandAbbreviations(tokens []Token) []string {
	// Common abbreviation expansions
	abbrevMap := map[string][]string{
		"st":   {"street", "saint"},
		"ave":  {"avenue"},
		"rd":   {"road"},
		"dr":   {"drive"},
		"blvd": {"boulevard"},
		"ln":   {"lane"},
		"ct":   {"court"},
		"pl":   {"place"},
		"sq":   {"square"},
		"apt":  {"apartment"},
		"n":    {"north"},
		"s":    {"south"},
		"e":    {"east"},
		"w":    {"west"},
		"ne":   {"northeast"},
		"nw":   {"northwest"},
		"se":   {"southeast"},
		"sw":   {"southwest"},
	}

	results := []string{""}

	for _, token := range tokens {
		if token.Type == TokenTypeWhitespace {
			for i := range results {
				results[i] += token.Value
			}
			continue
		}

		normalized := strings.ToLower(strings.TrimSpace(token.Value))
		expansions, found := abbrevMap[normalized]

		if found {
			// Create new results for each expansion
			var newResults []string
			for _, result := range results {
				// Keep original
				newResults = append(newResults, result+token.Value)
				// Add expansions
				for _, exp := range expansions {
					newResults = append(newResults, result+exp)
				}
			}
			results = newResults
		} else {
			// No expansion, add to all results
			for i := range results {
				results[i] += token.Value
			}
		}
	}

	return results
}

// removeFinalPeriods removes periods at the end of words
func removeFinalPeriods(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		words[i] = strings.TrimSuffix(word, ".")
	}
	return strings.Join(words, " ")
}

// ExpandAddressRoot is a convenience function using default options
func ExpandAddressRoot(input string) []string {
	return ExpandAddress(input, DefaultExpandOptions())
}
