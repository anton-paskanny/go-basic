# go-basic

## Overview
Small collection of basic Go CLI exercises/projects.

## Prerequisites
- Go installed (any recent stable version)

## Projects
- `1-converter`: simple converters (basic I/O and formatting)
- `2-calc`: a minimal calculator CLI
- `3-bin`: examples meant to be built as binaries

## Run
From each project directory:

```bash
cd 1-converter && go run .
cd 2-calc && go run .
cd 3-bin && go run .
```

## Build binaries
From a project directory:

```bash
go build -o ./build/app .
```

Then run:

```bash
./build/app
```