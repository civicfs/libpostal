package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/openvenues/libpostal/postal"
)

var (
	dataDir    = flag.String("datadir", "", "Custom data directory (default: system-specific)")
	downloadData = flag.Bool("download", false, "Download data files")
	parseMode  = flag.Bool("parse", false, "Parse mode: parse addresses from stdin")
	expandMode = flag.Bool("expand", false, "Expand mode: expand addresses from stdin")
	interactive = flag.Bool("interactive", false, "Interactive mode")
)

func main() {
	flag.Parse()

	// Download data if requested
	if *downloadData {
		fmt.Println("Downloading libpostal data files...")
		dm := postal.NewDataManager(*dataDir)
		if err := dm.Setup(); err != nil {
			log.Fatalf("Failed to download data: %v", err)
		}
		fmt.Println("Data download complete!")
		return
	}

	// Initialize library
	var err error
	if *dataDir != "" {
		err = postal.SetupWithDataDir(*dataDir)
	} else {
		err = postal.Setup()
	}

	if err != nil {
		log.Fatalf("Failed to initialize libpostal: %v\n", err)
	}
	defer postal.Teardown()

	// Parse mode
	if *parseMode {
		runParseMode()
		return
	}

	// Expand mode
	if *expandMode {
		runExpandMode()
		return
	}

	// Interactive mode
	if *interactive || flag.NArg() == 0 {
		runInteractive()
		return
	}

	// Parse command line argument
	address := strings.Join(flag.Args(), " ")
	parseAndPrint(address)
}

func runParseMode() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		address := scanner.Text()
		if address == "" {
			continue
		}
		parseAndPrint(address)
	}
}

func runExpandMode() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		address := scanner.Text()
		if address == "" {
			continue
		}
		expandAndPrint(address)
	}
}

func runInteractive() {
	fmt.Println("libpostal Interactive Mode")
	fmt.Println("Commands:")
	fmt.Println("  parse <address>  - Parse an address")
	fmt.Println("  expand <address> - Expand an address")
	fmt.Println("  quit            - Exit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 0 {
			continue
		}

		command := strings.ToLower(parts[0])

		switch command {
		case "quit", "exit", "q":
			return
		case "parse", "p":
			if len(parts) < 2 {
				fmt.Println("Usage: parse <address>")
				continue
			}
			parseAndPrint(parts[1])
		case "expand", "e":
			if len(parts) < 2 {
				fmt.Println("Usage: expand <address>")
				continue
			}
			expandAndPrint(parts[1])
		case "help", "h", "?":
			fmt.Println("Commands:")
			fmt.Println("  parse <address>  - Parse an address")
			fmt.Println("  expand <address> - Expand an address")
			fmt.Println("  quit            - Exit")
		default:
			// Try to parse as address
			parseAndPrint(line)
		}
	}
}

func parseAndPrint(address string) {
	response := postal.ParseAddress(address, postal.DefaultAddressParserOptions())

	fmt.Printf("Input: %s\n", address)
	fmt.Println("Components:")
	for _, comp := range response.Components {
		fmt.Printf("  %-15s: %s\n", comp.Label, comp.Value)
	}
	fmt.Println()
}

func expandAndPrint(address string) {
	expansions := postal.ExpandAddressRoot(address)

	fmt.Printf("Input: %s\n", address)
	fmt.Printf("Expansions (%d):\n", len(expansions))
	for i, exp := range expansions {
		fmt.Printf("  %d. %s\n", i+1, exp)
	}
	fmt.Println()
}
