# Bin CLI with JSONBin

This module provides a minimal CLI to create, fetch, and update JSON documents ("bins") using JSONBin v3, plus utilities for working with local JSON files.

## Prerequisites
- Go 1.20+
- JSONBin account and Master Key (v3)

## Setup
1. Create `.env` in `3-bin/` with your Master Key. Use single quotes to preserve `$`:
```bash
KEY='$2a$10$xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
```
2. Load env when running locally:
```bash
set -a && source .env && set +a
```

## CLI
Flags:
- `-create` Create a new bin from a JSON file. Optional: `-private` to make it private
- `-get`    Get a bin by id
- `-update` Update a bin by id from a JSON file
- `-id`     Bin id for get/update
- `-file`   Path to JSON file for create/update

Examples:
```bash
# Create a private bin from local JSON
go run . -create -file ./data/bins.json -private

# Fetch a bin by id
go run . -get -id <BIN_ID>

# Update a bin from local JSON
go run . -update -id <BIN_ID> -file ./data/bins.json
```

## Data
Sample payload lives at `data/bins.json`.
The CLI validates that files passed via `-file` are valid JSON before sending.

## Configuration
- Env var: `KEY` (required) – JSONBin Master Key
- API base: `https://api.jsonbin.io/v3`
- Headers: `X-Master-Key`, `X-Bin-Private`

## Project Structure
```
3-bin/
├── api/           # JSONBin client (POST/GET/PUT)
├── bins/          # Local data structures (example)
├── config/        # Env loader for KEY
├── data/          # Example JSON payloads
├── file/          # File helpers and JSON validation
├── storage/       # Local JSON persistence (example)
├── main.go        # CLI entrypoint
└── README.md
```

## Reference
- JSONBin overview: https://jsonbin.io/
