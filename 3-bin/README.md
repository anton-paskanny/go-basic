# Bin Storage System

This project demonstrates a simple bin storage system with JSON persistence and file validation capabilities.

## Features

### Storage Package (`storage/`)
- **SaveBins**: Saves bin list to JSON file with automatic directory creation
- **LoadBins**: Loads bin list from JSON file, returns empty list if file doesn't exist
- **IsJSONFile**: Checks if a file has a JSON extension
- **GetStoragePath**: Returns the current storage file path

### File Package (`file/`)
- **ReadAll**: Reads entire file contents
- **WriteAll**: Writes data to file with automatic directory creation
- **IsJSONFile**: Checks if a file has a JSON extension
- **ValidateJSONFile**: Validates that a file has .json extension and contains valid JSON
- **ReadJSONFile**: Reads and validates a JSON file

### Bins Package (`bins/`)
- **Bin**: Structure for storing bin information (ID, Name, Private, CreatedAt)
- **BinList**: Container for managing multiple bins
- **JSON Support**: All structures support JSON marshaling/unmarshaling

## Usage

### Running the Program
```bash
go run .
```

This will:
1. Create sample bins
2. Save them to `data/bins.json`
3. Load them back from the file
4. Display the results

### Example Output
The program creates a JSON file like:
```json
{
  "bins": [
    {
      "id": "bin001",
      "private": false,
      "created_at": "2025-08-30T09:11:21.468721+02:00",
      "name": "My first bin"
    }
  ]
}
```

### File Validation
```go
import "demo/bin/file"

// Check if file has JSON extension
if file.IsJSONFile("data.json") {
    // File is a JSON file
}

// Validate JSON file (checks extension and content)
err := file.ValidateJSONFile("data.json")
if err != nil {
    // File is not valid JSON
}
```

### Storage Operations
```go
import "demo/bin/storage"

// Create storage instance
storage := storage.New("my_bins.json")

// Save bins
err := storage.SaveBins(binList)

// Load bins
loadedBins, err := storage.LoadBins()
```

## Project Structure
```
3-bin/
├── bins/          # Bin data structures
├── storage/       # JSON persistence
├── file/          # File operations and validation
├── main.go        # Main program
└── README.md      # This file
```
