package main

import (
	"fmt"
	"log"

	"github.com/openvenues/libpostal/postal"
)

func main() {
	// Initialize the library
	if err := postal.Setup(); err != nil {
		log.Fatalf("Failed to setup postal: %v", err)
	}
	defer postal.Teardown()

	fmt.Println("=== libpostal Go Refactor Examples ===\n")

	// Example 1: Tokenization
	fmt.Println("1. Tokenization:")
	address := "123 Main Street, Apt 4B"
	tokens := postal.Tokenize(address, postal.DefaultTokenizeOptions())
	fmt.Printf("Input: %s\n", address)
	fmt.Printf("Tokens: ")
	for _, token := range tokens {
		fmt.Printf("[%s] ", token.Value)
	}
	fmt.Println("\n")

	// Example 2: Normalization
	fmt.Println("2. Normalization:")
	messy := "  Café  St.  Jean's  "
	normalized := postal.NormalizeString(messy)
	fmt.Printf("Input:      '%s'\n", messy)
	fmt.Printf("Normalized: '%s'\n\n", normalized)

	// Example 3: Address Expansion
	fmt.Println("3. Address Expansion:")
	input := "123 N. Main St."
	expansions := postal.ExpandAddressRoot(input)
	fmt.Printf("Input: %s\n", input)
	fmt.Println("Expansions:")
	for i, exp := range expansions {
		fmt.Printf("  %d. %s\n", i+1, exp)
	}
	fmt.Println()

	// Example 4: Address Parsing
	fmt.Println("4. Address Parsing:")
	fullAddress := "456 Park Avenue, Apartment 12, New York, NY 10022"
	parsed := postal.ParseAddress(fullAddress, postal.DefaultAddressParserOptions())
	fmt.Printf("Input: %s\n", fullAddress)
	fmt.Println("Parsed components:")
	for _, comp := range parsed.Components {
		fmt.Printf("  %-15s: %s\n", comp.Label, comp.Value)
	}
	fmt.Println()

	// Example 5: Duplicate Detection
	fmt.Println("5. Duplicate Detection:")
	name1 := "John Smith"
	name2 := "john smith"
	status, err := postal.IsNameDuplicate(name1, name2, postal.DefaultDuplicateOptions())
	if err != nil {
		log.Printf("Error checking duplicates: %v", err)
	} else {
		fmt.Printf("'%s' vs '%s': ", name1, name2)
		switch status {
		case postal.DuplicateStatusLikelyDuplicate:
			fmt.Println("Likely Duplicate ✓")
		case postal.DuplicateStatusPossibleDuplicate:
			fmt.Println("Possible Duplicate ~")
		case postal.DuplicateStatusNonDuplicate:
			fmt.Println("Not Duplicate ✗")
		}
	}
	fmt.Println()

	// Example 6: International Address
	fmt.Println("6. International Address (French):")
	french := "Quatre-vingt-douze Avenue des Champs-Élysées"
	frenchNorm := postal.NormalizeString(french)
	frenchTokens := postal.Tokenize(french, postal.DefaultTokenizeOptions())
	fmt.Printf("Original:   %s\n", french)
	fmt.Printf("Normalized: %s\n", frenchNorm)
	fmt.Printf("Tokens:     ")
	for _, token := range frenchTokens {
		fmt.Printf("[%s] ", token.Value)
	}
	fmt.Println("\n")

	// Example 7: CJK Address
	fmt.Println("7. CJK Address (Japanese):")
	japanese := "東京都渋谷区"
	japaneseTokens := postal.Tokenize(japanese, postal.DefaultTokenizeOptions())
	fmt.Printf("Original: %s\n", japanese)
	fmt.Printf("Tokens:   ")
	for _, token := range japaneseTokens {
		fmt.Printf("[%s:%d] ", token.Value, token.Type)
	}
	fmt.Println("\n")

	fmt.Println("=== All examples completed ===")
}
