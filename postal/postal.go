package postal

import (
	"errors"
	"sync"
)

// Package postal provides address parsing, normalization, and expansion
// for international addresses. This is a Go refactor of the libpostal C library.

var (
	setupMutex  sync.Mutex
	isSetup     bool
	dataManager *DataManager
)

// Setup initializes the postal library with the default data directory
func Setup() error {
	return SetupWithDataDir("")
}

// SetupWithDataDir initializes the postal library with a custom data directory
func SetupWithDataDir(dataDir string) error {
	setupMutex.Lock()
	defer setupMutex.Unlock()

	if isSetup {
		return nil
	}

	// Initialize data manager
	dataManager = NewDataManager(dataDir)

	// Download and extract data if needed
	if err := dataManager.Setup(); err != nil {
		return err
	}

	// Load data files
	if err := loadData(); err != nil {
		return err
	}

	isSetup = true
	return nil
}

// loadData loads the required data files
func loadData() error {
	// In a full implementation, this would:
	// - Load parser models from dataManager.GetDataPath("parser")
	// - Load language classifier from dataManager.GetDataPath("language_classifier")
	// - Load dictionaries and gazetteers
	// - Load transliteration rules
	// - Build internal data structures (tries, etc.)

	// For now, we can work without actual data files (using built-in rules)
	// In production, you would download and load the actual data

	return nil
}

// SetupParser initializes only the parser component
func SetupParser() error {
	return Setup()
}

// SetupLanguageClassifier initializes only the language classifier
func SetupLanguageClassifier() error {
	return Setup()
}

// Teardown cleans up resources used by the library
func Teardown() {
	setupMutex.Lock()
	defer setupMutex.Unlock()

	isSetup = false
}

// IsSetup returns whether the library has been initialized
func IsSetup() bool {
	setupMutex.Lock()
	defer setupMutex.Unlock()
	return isSetup
}

// DuplicateStatus represents the result of a duplicate check
type DuplicateStatus int

const (
	DuplicateStatusNonDuplicate DuplicateStatus = iota
	DuplicateStatusLikelyDuplicate
	DuplicateStatusPossibleDuplicate
)

// DuplicateOptions configures duplicate detection
type DuplicateOptions struct {
	Languages []string
}

// DefaultDuplicateOptions returns default duplicate detection options
func DefaultDuplicateOptions() DuplicateOptions {
	return DuplicateOptions{
		Languages: nil,
	}
}

// IsNameDuplicate checks if two names are likely duplicates
func IsNameDuplicate(name1, name2 string, options DuplicateOptions) (DuplicateStatus, error) {
	if !isSetup {
		return DuplicateStatusNonDuplicate, errors.New("postal library not initialized - call Setup() first")
	}

	// Normalize both names
	norm1 := NormalizeString(name1)
	norm2 := NormalizeString(name2)

	// Exact match after normalization
	if norm1 == norm2 {
		return DuplicateStatusLikelyDuplicate, nil
	}

	// Calculate similarity
	similarity := calculateSimilarity(norm1, norm2)

	if similarity > 0.9 {
		return DuplicateStatusLikelyDuplicate, nil
	} else if similarity > 0.7 {
		return DuplicateStatusPossibleDuplicate, nil
	}

	return DuplicateStatusNonDuplicate, nil
}

// calculateSimilarity computes a simple similarity score between two strings
func calculateSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	// Simple Jaccard similarity on character bigrams
	bigrams1 := getBigrams(s1)
	bigrams2 := getBigrams(s2)

	intersection := 0
	for bg := range bigrams1 {
		if bigrams2[bg] {
			intersection++
		}
	}

	union := len(bigrams1) + len(bigrams2) - intersection
	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

// getBigrams extracts character bigrams from a string
func getBigrams(s string) map[string]bool {
	bigrams := make(map[string]bool)
	runes := []rune(s)

	for i := 0; i < len(runes)-1; i++ {
		bigram := string(runes[i:i+2])
		bigrams[bigram] = true
	}

	return bigrams
}

// NearDupeHashOptions configures near-duplicate hashing
type NearDupeHashOptions struct {
	Languages []string
}

// DefaultNearDupeHashOptions returns default near-dupe hash options
func DefaultNearDupeHashOptions() NearDupeHashOptions {
	return NearDupeHashOptions{
		Languages: nil,
	}
}

// NearDupeHashes generates hashes for near-duplicate detection
func NearDupeHashes(labels []string, values []string, options NearDupeHashOptions) ([]string, error) {
	if !isSetup {
		return nil, errors.New("postal library not initialized - call Setup() first")
	}

	if len(labels) != len(values) {
		return nil, errors.New("labels and values must have the same length")
	}

	var hashes []string

	// Generate hashes based on normalized components
	for i, value := range values {
		if value == "" {
			continue
		}

		normalized := NormalizeString(value)
		hash := labels[i] + ":" + normalized
		hashes = append(hashes, hash)
	}

	return hashes, nil
}
