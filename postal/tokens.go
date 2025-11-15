package postal

import (
	"unicode"
)

// TokenType represents the type of token
type TokenType int

const (
	TokenTypeWord TokenType = iota
	TokenTypeAbbreviation
	TokenTypePunctuation
	TokenTypeWhitespace
	TokenTypeNumeric
	TokenTypeIdeographic
	TokenTypeUnknown
)

// Token represents a parsed token from an address string
type Token struct {
	Value  string
	Type   TokenType
	Offset int
	Length int
}

// TokenizeOptions configures tokenization behavior
type TokenizeOptions struct {
	KeepWhitespace  bool
	KeepPunctuation bool
}

// DefaultTokenizeOptions returns default tokenization options
func DefaultTokenizeOptions() TokenizeOptions {
	return TokenizeOptions{
		KeepWhitespace:  false,
		KeepPunctuation: true,
	}
}

// Tokenize splits an input string into tokens
func Tokenize(input string, options TokenizeOptions) []Token {
	if input == "" {
		return nil
	}

	var tokens []Token
	currentToken := ""
	currentType := TokenTypeUnknown
	offset := 0

	for i, r := range input {
		var tokenType TokenType

		switch {
		case unicode.IsSpace(r):
			tokenType = TokenTypeWhitespace
		case unicode.IsPunct(r):
			tokenType = TokenTypePunctuation
		case unicode.IsDigit(r):
			tokenType = TokenTypeNumeric
		case unicode.IsLetter(r):
			// Check if ideographic (CJK, etc.)
			if isIdeographic(r) {
				tokenType = TokenTypeIdeographic
			} else {
				tokenType = TokenTypeWord
			}
		default:
			tokenType = TokenTypeUnknown
		}

		// If token type changed or ideographic (each gets its own token)
		if tokenType != currentType || tokenType == TokenTypeIdeographic {
			if currentToken != "" {
				if shouldKeepToken(currentType, options) {
					tokens = append(tokens, Token{
						Value:  currentToken,
						Type:   currentType,
						Offset: offset,
						Length: len(currentToken),
					})
				}
			}
			currentToken = string(r)
			currentType = tokenType
			offset = i
		} else {
			currentToken += string(r)
		}
	}

	// Add final token
	if currentToken != "" && shouldKeepToken(currentType, options) {
		tokens = append(tokens, Token{
			Value:  currentToken,
			Type:   currentType,
			Offset: offset,
			Length: len(currentToken),
		})
	}

	return tokens
}

func shouldKeepToken(tokenType TokenType, options TokenizeOptions) bool {
	if !options.KeepWhitespace && tokenType == TokenTypeWhitespace {
		return false
	}
	if !options.KeepPunctuation && tokenType == TokenTypePunctuation {
		return false
	}
	return true
}

func isIdeographic(r rune) bool {
	// CJK Unified Ideographs ranges
	return (r >= 0x4E00 && r <= 0x9FFF) || // CJK Unified Ideographs
		(r >= 0x3400 && r <= 0x4DBF) || // CJK Unified Ideographs Extension A
		(r >= 0x20000 && r <= 0x2A6DF) || // CJK Unified Ideographs Extension B
		(r >= 0x2A700 && r <= 0x2B73F) || // CJK Unified Ideographs Extension C
		(r >= 0x2B740 && r <= 0x2B81F) || // CJK Unified Ideographs Extension D
		(r >= 0x2B820 && r <= 0x2CEAF) || // CJK Unified Ideographs Extension E
		(r >= 0xF900 && r <= 0xFAFF) || // CJK Compatibility Ideographs
		(r >= 0x2F800 && r <= 0x2FA1F) // CJK Compatibility Ideographs Supplement
}

// TokensToString converts a slice of tokens back to a string
func TokensToString(tokens []Token) string {
	result := ""
	for _, token := range tokens {
		result += token.Value
	}
	return result
}
