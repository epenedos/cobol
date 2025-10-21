package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// AddressRecord represents a single address entry
type AddressRecord struct {
	LastName  string
	FirstName string
	Street    string
	City      string
	State     string
	Zip       string
}

// Config holds file paths
type Config struct {
	InputFile  string
	OutputFile string
}

// Default configuration matching COBOL program paths
var defaultConfig = Config{
	InputFile:  "/nfs_dir/input/info.csv",
	OutputFile: "/nfs_dir/output/output.txt",
}

// formatFixedWidth formats an address record into a 160-character fixed-width string
// matching the COBOL output format:
// Last Name (25) + Space (5) + First Name (15) + Space (5) + Street (30) +
// Space (5) + City (15) + Space (5) + State (3) + Space (5) + Zip (10) + Space (38)
func formatFixedWidth(record AddressRecord) string {
	var builder strings.Builder
	builder.Grow(160) // Pre-allocate for efficiency

	// Last name - 25 characters
	builder.WriteString(padRight(record.LastName, 25))
	// Filler - 5 spaces
	builder.WriteString(padRight("", 5))
	// First name - 15 characters
	builder.WriteString(padRight(record.FirstName, 15))
	// Filler - 5 spaces
	builder.WriteString(padRight("", 5))
	// Street - 30 characters
	builder.WriteString(padRight(record.Street, 30))
	// Filler - 5 spaces
	builder.WriteString(padRight("", 5))
	// City - 15 characters
	builder.WriteString(padRight(record.City, 15))
	// Filler - 5 spaces
	builder.WriteString(padRight("", 5))
	// State - 3 characters
	builder.WriteString(padRight(record.State, 3))
	// Filler - 5 spaces
	builder.WriteString(padRight("", 5))
	// Zip - 10 characters
	builder.WriteString(padRight(record.Zip, 10))
	// Filler - 38 spaces (total = 160 characters)
	builder.WriteString(padRight("", 38))

	return builder.String()
}

// padRight pads a string with spaces to the right to reach the specified length
// If the string is longer, it truncates to fit
func padRight(s string, length int) string {
	if len(s) >= length {
		return s[:length]
	}
	return s + strings.Repeat(" ", length-len(s))
}

// processCSV reads the CSV file and writes formatted fixed-width output
func processCSV(config Config) error {
	// Open input file
	inputFile, err := os.Open(config.InputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	// Create output file
	outputFile, err := os.Create(config.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Create CSV reader
	csvReader := csv.NewReader(inputFile)
	csvReader.FieldsPerRecord = 6 // Expect 6 fields per record

	recordCount := 0

	// Read and process each record
	for {
		fields, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV at record %d: %w", recordCount+1, err)
		}

		// Ensure we have exactly 6 fields
		if len(fields) != 6 {
			log.Printf("Warning: record %d has %d fields (expected 6), skipping", recordCount+1, len(fields))
			continue
		}

		// Create address record
		record := AddressRecord{
			LastName:  strings.TrimSpace(fields[0]),
			FirstName: strings.TrimSpace(fields[1]),
			Street:    strings.TrimSpace(fields[2]),
			City:      strings.TrimSpace(fields[3]),
			State:     strings.TrimSpace(fields[4]),
			Zip:       strings.TrimSpace(fields[5]),
		}

		// Format and write the record
		formattedLine := formatFixedWidth(record)
		if _, err := outputFile.WriteString(formattedLine + "\n"); err != nil {
			return fmt.Errorf("error writing record %d: %w", recordCount+1, err)
		}

		recordCount++
	}

	log.Printf("Successfully processed %d records", recordCount)
	return nil
}

func main() {
	// Allow override via command-line arguments
	config := defaultConfig
	if len(os.Args) > 2 {
		config.InputFile = os.Args[1]
		config.OutputFile = os.Args[2]
	}

	log.Printf("Reading from: %s", config.InputFile)
	log.Printf("Writing to: %s", config.OutputFile)

	if err := processCSV(config); err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Println("Processing complete")
}
