# COBOL to Go Reengineering: CSV Report Generator

This repository contains both the original COBOL program and its modern Go reengineering.

## Overview

Both programs perform the same function: read a CSV file containing address information and convert it to a fixed-width text format suitable for legacy systems or reports.

## Files

- `report.cbl` - Original COBOL program
- `report.go` - Reengineered Go implementation
- `test_data/sample.csv` - Sample input data
- `test_data/output.txt` - Example output

## Functionality

### Input Format (CSV)
```
LastName,FirstName,Street,City,State,Zip
```

### Output Format (Fixed-Width, 160 characters per line)
- Last Name: 25 characters
- Spacing: 5 characters
- First Name: 15 characters
- Spacing: 5 characters
- Street: 30 characters
- Spacing: 5 characters
- City: 15 characters
- Spacing: 5 characters
- State: 3 characters
- Spacing: 5 characters
- Zip: 10 characters
- Trailing spaces: 38 characters

**Total: 160 characters per record**

## Key Differences: COBOL vs Go

### 1. File Path Configuration

**COBOL (Hardcoded):**
```cobol
SELECT INPUT-FILE
    ASSIGN TO "/nfs_dir/input/info.csv"
```

**Go (Flexible):**
```go
// Default paths with command-line override support
go run report.go [input.csv] [output.txt]
```

### 2. Data Structures

**COBOL (Fixed-length fields):**
```cobol
01  SEPARATE-IT.
    05 LAST_NAME        PIC X(25).
    05 FIRST_NAME       PIC X(15).
```

**Go (Struct with dynamic sizing):**
```go
type AddressRecord struct {
    LastName  string
    FirstName string
    // ...
}
```

### 3. CSV Parsing

**COBOL (Manual parsing):**
```cobol
UNSTRING INPUT-RECORD DELIMITED BY ","
   INTO LAST_NAME, FIRST_NAME, STREET_ADDR,
   CITY, STATE, ZIP.
```

**Go (Standard library):**
```go
csvReader := csv.NewReader(inputFile)
fields, err := csvReader.Read()
```

### 4. Error Handling

**COBOL (Minimal):**
```cobol
READ INPUT-FILE AT END GO TO END-ROUTINE.
```

**Go (Comprehensive):**
```go
if err != nil {
    return fmt.Errorf("error reading CSV: %w", err)
}
```

### 5. Control Flow

**COBOL (GO TO statements):**
```cobol
READ-ROUTINE.
    READ INPUT-FILE AT END GO TO END-ROUTINE.
    ...
    GO TO READ-ROUTINE.
```

**Go (Modern loops):**
```go
for {
    fields, err := csvReader.Read()
    if err == io.EOF {
        break
    }
    // Process record
}
```

### 6. String Formatting

**COBOL (Implicit padding via PIC clauses):**
```cobol
05 OUT-LAST-NAME     PIC X(25).
```

**Go (Explicit padding function):**
```go
func padRight(s string, length int) string {
    if len(s) >= length {
        return s[:length]
    }
    return s + strings.Repeat(" ", length-len(s))
}
```

## Running the Programs

### COBOL
```bash
# Requires GnuCOBOL or similar compiler
cobc -x report.cbl
./report
```

### Go
```bash
# Using default paths
go run report.go

# Using custom paths
go run report.go test_data/sample.csv test_data/output.txt

# Build executable
go build -o report report.go
./report test_data/sample.csv test_data/output.txt
```

## Advantages of the Go Implementation

1. **Modern Error Handling**: Detailed error messages with context
2. **Flexibility**: Command-line arguments for file paths
3. **Type Safety**: Compile-time checks for data structures
4. **Standard Library**: Built-in CSV parsing
5. **Logging**: Progress reporting and diagnostics
6. **Maintainability**: Clear, readable code structure
7. **Performance**: Efficient string building and memory management
8. **Portability**: Cross-platform compatibility
9. **Testing**: Easy to write unit tests
10. **Deployment**: Single binary with no runtime dependencies

## Functional Equivalence

Both programs produce identical output for the same input data. The Go version adds:
- Input validation (field count checking)
- Detailed logging
- Better error messages
- Flexible configuration

## Testing

```bash
# Test with sample data
go run report.go test_data/sample.csv test_data/output.txt

# Verify output format
cat test_data/output.txt
```

## Example Output

```
Smith                         John                123 Main Street                    Springfield         IL      62701
Doe                           Jane                456 Oak Avenue                     Chicago             IL      60601
```

Each line is exactly 160 characters with proper field alignment.

## Architecture Comparison

| Aspect | COBOL | Go |
|--------|-------|-----|
| Paradigm | Procedural | Procedural + OO |
| Memory | Fixed allocation | Dynamic allocation |
| Error Handling | AT END clause | Error values + logging |
| String Processing | MOVE + PIC | strings package |
| CSV Parsing | UNSTRING | encoding/csv |
| Control Flow | GO TO | for/if/break |
| Modularity | Divisions | Functions + packages |
| Testing | External | Built-in testing |

## Conclusion

This reengineering demonstrates how legacy COBOL business logic can be modernized using Go while maintaining functional equivalence. The Go version provides better maintainability, error handling, and flexibility while preserving the exact output format required by downstream systems.
