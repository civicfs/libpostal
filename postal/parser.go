package postal

import (
	"strings"
)

// AddressComponent represents a component label for parsed addresses
type AddressComponent string

const (
	ComponentHouseNumber AddressComponent = "house_number"
	ComponentRoad        AddressComponent = "road"
	ComponentUnit        AddressComponent = "unit"
	ComponentLevel       AddressComponent = "level"
	ComponentStaircase   AddressComponent = "staircase"
	ComponentEntrance    AddressComponent = "entrance"
	ComponentPoBox       AddressComponent = "po_box"
	ComponentPostcode    AddressComponent = "postcode"
	ComponentSuburb      AddressComponent = "suburb"
	ComponentCityDistrict AddressComponent = "city_district"
	ComponentCity        AddressComponent = "city"
	ComponentIsland      AddressComponent = "island"
	ComponentStateDistrict AddressComponent = "state_district"
	ComponentState       AddressComponent = "state"
	ComponentCountryRegion AddressComponent = "country_region"
	ComponentCountry     AddressComponent = "country"
	ComponentWorldRegion AddressComponent = "world_region"
)

// ParsedComponent represents a single parsed component of an address
type ParsedComponent struct {
	Label string
	Value string
}

// AddressParserResponse contains the parsed address components
type AddressParserResponse struct {
	Components []ParsedComponent
}

// AddressParserOptions configures address parsing behavior
type AddressParserOptions struct {
	Language string
	Country  string
}

// DefaultAddressParserOptions returns default parser options
func DefaultAddressParserOptions() AddressParserOptions {
	return AddressParserOptions{
		Language: "", // Auto-detect
		Country:  "", // Auto-detect
	}
}

// ParseAddress parses an address string into labeled components
// This is a simplified implementation - the full C library uses ML models
func ParseAddress(input string, options AddressParserOptions) *AddressParserResponse {
	if input == "" {
		return &AddressParserResponse{Components: nil}
	}

	// Normalize the input
	normalized := NormalizeString(input)

	// Tokenize
	tokens := Tokenize(normalized, TokenizeOptions{
		KeepWhitespace:  true,
		KeepPunctuation: true,
	})

	// Simple rule-based parsing (in real implementation, would use ML model)
	components := parseAddressComponents(tokens, normalized)

	return &AddressParserResponse{
		Components: components,
	}
}

// parseAddressComponents performs simple rule-based component extraction
func parseAddressComponents(tokens []Token, normalized string) []ParsedComponent {
	var components []ParsedComponent

	// Split by comma for basic structure
	parts := strings.Split(normalized, ",")

	if len(parts) == 0 {
		return components
	}

	// First part typically contains street address
	if len(parts) > 0 {
		streetPart := strings.TrimSpace(parts[0])
		parsed := parseStreetAddress(streetPart)
		components = append(components, parsed...)
	}

	// Middle parts often contain city/locality
	if len(parts) > 1 {
		for i := 1; i < len(parts)-1; i++ {
			part := strings.TrimSpace(parts[i])
			if part != "" {
				// Could be suburb or city
				components = append(components, ParsedComponent{
					Label: string(ComponentCity),
					Value: part,
				})
			}
		}
	}

	// Last part often contains state and postcode
	if len(parts) > 1 {
		lastPart := strings.TrimSpace(parts[len(parts)-1])
		parsed := parseStatePostcode(lastPart)
		components = append(components, parsed...)
	}

	return components
}

// parseStreetAddress extracts house number and road from street address
func parseStreetAddress(street string) []ParsedComponent {
	var components []ParsedComponent

	tokens := strings.Fields(street)
	if len(tokens) == 0 {
		return components
	}

	// Check if first token is a number (house number)
	if isNumeric(tokens[0]) {
		components = append(components, ParsedComponent{
			Label: string(ComponentHouseNumber),
			Value: tokens[0],
		})

		// Rest is the road
		if len(tokens) > 1 {
			road := strings.Join(tokens[1:], " ")
			components = append(components, ParsedComponent{
				Label: string(ComponentRoad),
				Value: road,
			})
		}
	} else {
		// Entire string is the road
		components = append(components, ParsedComponent{
			Label: string(ComponentRoad),
			Value: street,
		})
	}

	return components
}

// parseStatePostcode extracts state and postcode from the last part
func parseStatePostcode(part string) []ParsedComponent {
	var components []ParsedComponent

	tokens := strings.Fields(part)
	if len(tokens) == 0 {
		return components
	}

	// Look for postcode (numeric or alphanumeric pattern)
	for i := len(tokens) - 1; i >= 0; i-- {
		if looksLikePostcode(tokens[i]) {
			// Found postcode
			components = append(components, ParsedComponent{
				Label: string(ComponentPostcode),
				Value: tokens[i],
			})

			// Everything before is state/region
			if i > 0 {
				state := strings.Join(tokens[:i], " ")
				components = append([]ParsedComponent{{
					Label: string(ComponentState),
					Value: state,
				}}, components...)
			}
			return components
		}
	}

	// No postcode found, treat as state
	components = append(components, ParsedComponent{
		Label: string(ComponentState),
		Value: part,
	})

	return components
}

// isNumeric checks if a string contains only digits
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// looksLikePostcode checks if a string looks like a postcode
func looksLikePostcode(s string) bool {
	if len(s) < 3 {
		return false
	}

	// Check for common postcode patterns
	digitCount := 0
	for _, r := range s {
		if r >= '0' && r <= '9' {
			digitCount++
		}
	}

	// Most postcodes have at least 2 digits
	return digitCount >= 2
}

// GetComponent extracts a specific component value from the response
func (r *AddressParserResponse) GetComponent(label string) string {
	for _, comp := range r.Components {
		if comp.Label == label {
			return comp.Value
		}
	}
	return ""
}

// GetComponents extracts all values for a specific component label
func (r *AddressParserResponse) GetComponents(label string) []string {
	var values []string
	for _, comp := range r.Components {
		if comp.Label == label {
			values = append(values, comp.Value)
		}
	}
	return values
}
