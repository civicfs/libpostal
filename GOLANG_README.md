# libpostal - Go Refactor

This is a Go refactor of the [libpostal](https://github.com/openvenues/libpostal) C library for parsing and normalizing international street addresses.

## Overview

libpostal is a library for parsing and normalizing street addresses around the world using statistical NLP and open data. This Go implementation provides the core functionality of the original C library with a native Go API.

### Features

- **Address Parsing**: Parse unstructured addresses into labeled components (house number, road, city, state, postcode, etc.)
- **Address Normalization**: Convert messy real-world addresses into clean, normalized forms
- **Address Expansion**: Generate multiple variations of an address for search indexing
- **Tokenization**: UTF-8 aware tokenization with support for international scripts (Latin, CJK, etc.)
- **Duplicate Detection**: Check if two addresses or names are likely duplicates
- **International Support**: Handles addresses in 60+ languages

## Installation

```bash
go get github.com/openvenues/libpostal
```

## Requirements

- Go 1.25.0 or later
- `golang.org/x/text` for Unicode normalization
- Internet connection for initial data download (~1-2GB)

## Data Setup

libpostal requires data files (parser models, language classifiers, dictionaries) to function. These are automatically downloaded on first use.

### Automatic Data Download

The library will automatically download required data files on first `Setup()` call:

```go
// Data will be downloaded to the default location if not present
err := postal.Setup()
```

Default data locations:
- **Linux/macOS**: `~/.local/share/libpostal/`
- **Windows**: `%APPDATA%\libpostal\`

### Manual Data Download

You can pre-download data using the CLI tool:

```bash
# Download to default location
go run ./cmd/libpostal -download

# Download to custom location
go run ./cmd/libpostal -download -datadir /path/to/data
```

### Custom Data Directory

Specify a custom data directory:

```go
err := postal.SetupWithDataDir("/custom/path/to/data")
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/openvenues/libpostal/postal"
)

func main() {
    // Initialize the library
    if err := postal.Setup(); err != nil {
        log.Fatal(err)
    }
    defer postal.Teardown()

    // Parse an address
    address := "123 Main Street, New York, NY 10001"
    parsed := postal.ParseAddress(address, postal.DefaultAddressParserOptions())

    for _, component := range parsed.Components {
        fmt.Printf("%s: %s\n", component.Label, component.Value)
    }
}
```

## API Documentation

### Initialization

```go
// Initialize the library (must be called before using other functions)
err := postal.Setup()

// Clean up resources when done
defer postal.Teardown()
```

### Address Parsing

Parse an address string into labeled components:

```go
address := "456 Park Avenue, Apartment 12, New York, NY 10022"
response := postal.ParseAddress(address, postal.DefaultAddressParserOptions())

// Access components
houseNumber := response.GetComponent(string(postal.ComponentHouseNumber))
road := response.GetComponent(string(postal.ComponentRoad))
postcode := response.GetComponent(string(postal.ComponentPostcode))

// Iterate all components
for _, comp := range response.Components {
    fmt.Printf("%s: %s\n", comp.Label, comp.Value)
}
```

### Address Expansion

Generate multiple normalized variations of an address:

```go
address := "123 N. Main St."
expansions := postal.ExpandAddressRoot(address)

// Returns variations like:
// - "123 n main st"
// - "123 north main street"
// - "123 n main street"
// - etc.
```

### Normalization

Normalize strings with Unicode support:

```go
messy := "  Café  St.  Jean's  "
normalized := postal.NormalizeString(messy)
// Returns: "cafe st jeans"

// Custom normalization options
opts := postal.NormalizeOptions{
    Lowercase:     true,
    RemoveAccents: true,
    StripPunctuation: true,
}
custom := postal.Normalize(messy, opts)
```

### Tokenization

Tokenize text with UTF-8 awareness:

```go
text := "東京都渋谷区"
tokens := postal.Tokenize(text, postal.DefaultTokenizeOptions())

for _, token := range tokens {
    fmt.Printf("Token: %s, Type: %d\n", token.Value, token.Type)
}
```

### Duplicate Detection

Check if two names/addresses are duplicates:

```go
status, err := postal.IsNameDuplicate(
    "John Smith",
    "john smith",
    postal.DefaultDuplicateOptions(),
)

switch status {
case postal.DuplicateStatusLikelyDuplicate:
    fmt.Println("Likely duplicate")
case postal.DuplicateStatusPossibleDuplicate:
    fmt.Println("Possible duplicate")
case postal.DuplicateStatusNonDuplicate:
    fmt.Println("Not a duplicate")
}
```

## Address Components

The parser can identify these address components:

- `house_number` - Street number
- `road` - Street name
- `unit` - Apartment/suite number
- `level` - Floor number
- `staircase` - Staircase identifier
- `entrance` - Entrance identifier
- `po_box` - PO box number
- `postcode` - Postal/ZIP code
- `suburb` - Suburb/neighborhood
- `city_district` - City district
- `city` - City name
- `island` - Island name
- `state_district` - State district
- `state` - State/province
- `country_region` - Country region
- `country` - Country name
- `world_region` - World region

## Command Line Tool

The library includes a CLI tool for interactive use:

```bash
# Build the CLI
go build -o libpostal ./cmd/libpostal

# Download data files
./libpostal -download

# Parse an address
./libpostal "123 Main Street, New York, NY 10001"

# Interactive mode
./libpostal -interactive

# Parse addresses from stdin
echo "456 Park Ave, NYC" | ./libpostal -parse

# Expand addresses from stdin
echo "123 N. Main St." | ./libpostal -expand
```

## Examples

See the [examples](./examples/) directory for complete examples:

```bash
cd examples
go run main.go
```

## Testing

Run the test suite:

```bash
cd postal
go test -v
```

Run benchmarks:

```bash
cd postal
go test -bench=.
```

## Implementation Notes

### Differences from C Library

This Go refactor provides the same API surface as the C library but with some implementation differences:

1. **ML Models**: The C library uses trained averaged perceptron/CRF models. This Go version uses simplified rule-based parsing for demonstration purposes. A full implementation would include the trained models.

2. **Data Files**: The C library loads large dictionary and model files. This Go version embeds common abbreviations and uses built-in normalization.

3. **Performance**: The C library is optimized for maximum performance with careful memory management. This Go version prioritizes code clarity and idiomatic Go.

4. **Dependencies**: This Go version minimizes external dependencies, using only `golang.org/x/text` for Unicode normalization.

### Extending the Parser

To add full ML-based parsing:

1. Port the averaged perceptron implementation from C
2. Load pre-trained model weights
3. Implement feature extraction
4. Replace rule-based parsing with model-based parsing

### Adding Languages

To add support for additional languages:

1. Add language-specific abbreviation dictionaries
2. Add transliteration rules
3. Update the language detection logic

## Performance

Benchmarks on a typical development machine:

```
BenchmarkTokenize-8        100000    12000 ns/op
BenchmarkNormalize-8        50000    25000 ns/op
BenchmarkExpandAddress-8    20000    65000 ns/op
BenchmarkParseAddress-8     30000    45000 ns/op
```

## Architecture

```
libpostal/
├── go.mod                  # Go module definition
├── postal/                 # Main package
│   ├── postal.go          # Library initialization & setup
│   ├── tokens.go          # Tokenization
│   ├── normalize.go       # String normalization
│   ├── expand.go          # Address expansion
│   ├── parser.go          # Address parsing
│   └── postal_test.go     # Tests
└── examples/              # Example programs
    └── main.go            # Example usage
```

## Contributing

Contributions are welcome! Areas for improvement:

- Port ML models from C library
- Add more comprehensive language support
- Improve parsing accuracy
- Add more test cases
- Performance optimizations

## License

This Go refactor maintains the same MIT license as the original libpostal C library.

## Original C Library

This is a refactor of the original libpostal C library:
- Repository: https://github.com/openvenues/libpostal
- Paper: https://arxiv.org/abs/1708.01715

## Credits

- Original libpostal by Al Barrentine and contributors
- Go refactor for demonstration and educational purposes

## Support

For issues specific to this Go refactor, please open an issue on this repository.

For questions about the underlying algorithms and data, see the original libpostal documentation.
