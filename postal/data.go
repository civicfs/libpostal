package postal

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const (
	// Default data URL - libpostal data hosted on GitHub releases
	DefaultDataURL     = "https://github.com/openvenues/libpostal/releases/download"
	DataVersion        = "v1.0.0"
	BaseDataFile       = "libpostal_data.tar.gz"
	ParserDataFile     = "parser.tar.gz"
	LanguageDataFile   = "language_classifier.tar.gz"

	// Default data directory
	DefaultDataDir = "libpostal_data"
)

// DataManager handles downloading and managing libpostal data files
type DataManager struct {
	DataDir string
	BaseURL string
}

// NewDataManager creates a new data manager
func NewDataManager(dataDir string) *DataManager {
	if dataDir == "" {
		// Use system-appropriate data directory
		dataDir = getDefaultDataDir()
	}

	return &DataManager{
		DataDir: dataDir,
		BaseURL: DefaultDataURL,
	}
}

// getDefaultDataDir returns the default data directory for the OS
func getDefaultDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	switch runtime.GOOS {
	case "linux", "darwin":
		return filepath.Join(homeDir, ".local", "share", "libpostal")
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "libpostal")
		}
		return filepath.Join(homeDir, "libpostal_data")
	default:
		return filepath.Join(homeDir, "libpostal_data")
	}
}

// Setup downloads and extracts required data files
func (dm *DataManager) Setup() error {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dm.DataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Check if data already exists
	if dm.DataExists() {
		fmt.Println("Data files found, skipping download.")
		return nil
	}

	fmt.Println("Downloading libpostal data files (this may take a while)...")
	fmt.Println("Note: Data files are large (~1-2GB total)")

	// Download base data
	if err := dm.DownloadAndExtract(BaseDataFile); err != nil {
		fmt.Printf("Warning: Failed to download base data: %v\n", err)
		fmt.Println("Continuing with built-in rules only...")
	}

	// Download parser data
	if err := dm.DownloadAndExtract(ParserDataFile); err != nil {
		fmt.Printf("Warning: Failed to download parser data: %v\n", err)
		fmt.Println("Continuing with built-in rules only...")
	}

	// Download language classifier data
	if err := dm.DownloadAndExtract(LanguageDataFile); err != nil {
		fmt.Printf("Warning: Failed to download language classifier data: %v\n", err)
		fmt.Println("Continuing with built-in rules only...")
	}

	fmt.Println("Data download complete!")

	return nil
}

// DataExists checks if required data files exist
func (dm *DataManager) DataExists() bool {
	// Check for key directories that should exist after extraction
	basePaths := []string{
		filepath.Join(dm.DataDir, "address_expansions"),
		filepath.Join(dm.DataDir, "numex"),
		filepath.Join(dm.DataDir, "transliteration"),
	}

	parserPath := filepath.Join(dm.DataDir, "address_parser")
	languagePath := filepath.Join(dm.DataDir, "language_classifier")

	// Check base data
	baseExists := true
	for _, path := range basePaths {
		if info, err := os.Stat(path); err != nil || !info.IsDir() {
			baseExists = false
			break
		}
	}

	// Check parser
	parserExists := false
	if info, err := os.Stat(parserPath); err == nil && info.IsDir() {
		parserExists = true
	}

	// Check language classifier
	languageExists := false
	if info, err := os.Stat(languagePath); err == nil && info.IsDir() {
		languageExists = true
	}

	return baseExists && parserExists && languageExists
}

// DownloadAndExtract downloads and extracts a data file
func (dm *DataManager) DownloadAndExtract(filename string) error {
	url := fmt.Sprintf("%s/%s/%s", dm.BaseURL, DataVersion, filename)
	tempFile := filepath.Join(dm.DataDir, filename)

	// Download the file
	fmt.Printf("Downloading %s...\n", filename)
	if err := dm.downloadFile(url, tempFile); err != nil {
		return err
	}

	// Extract the file
	fmt.Printf("Extracting %s...\n", filename)
	if err := dm.extractTarGz(tempFile, dm.DataDir); err != nil {
		return err
	}

	// Clean up the archive
	os.Remove(tempFile)

	return nil
}

// downloadFile downloads a file from a URL
func (dm *DataManager) downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Create a progress reader
	size := resp.ContentLength
	var downloaded int64

	buf := make([]byte, 32*1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, writeErr := out.Write(buf[:n])
			if writeErr != nil {
				return writeErr
			}
			downloaded += int64(n)

			// Print progress
			if size > 0 {
				percent := float64(downloaded) / float64(size) * 100
				fmt.Printf("\rProgress: %.1f%%", percent)
			}
		}
		if err != nil {
			if err == io.EOF {
				fmt.Println() // New line after progress
				break
			}
			return err
		}
	}

	return nil
}

// extractTarGz extracts a .tar.gz file
func (dm *DataManager) extractTarGz(tarGzPath, destDir string) error {
	file, err := os.Open(tarGzPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			// Create parent directories
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}

			outFile, err := os.Create(target)
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}

	return nil
}

// GetDataPath returns the path to a data file
func (dm *DataManager) GetDataPath(subpath string) string {
	return filepath.Join(dm.DataDir, subpath)
}

// SetupDataDir is a convenience function to set up data with a custom directory
func SetupDataDir(dataDir string) error {
	dm := NewDataManager(dataDir)
	return dm.Setup()
}

// SetupDefaultData sets up data in the default directory
func SetupDefaultData() error {
	dm := NewDataManager("")
	return dm.Setup()
}
