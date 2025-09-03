package main

import (
	"demo/bin/api"
	b "demo/bin/bins"
	"demo/bin/config"
	f "demo/bin/file"
	s "demo/bin/storage"
	"fmt"
	"log"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create API client with configuration
	apiClient := api.New(cfg)
	_ = apiClient // TODO: use API client when methods are implemented

	// Create a new list
	binList := b.NewList()

	// Create several example bins
	bin1 := b.NewBin("bin001", "My first bin", false)
	bin2 := b.NewBin("bin002", "Private bin", true)
	bin3 := b.NewBin("bin003", "Work bin", false)

	// Add bins to the list
	binList.Add(*bin1)
	binList.Add(*bin2)
	binList.Add(*bin3)

	// Print all bins
	fmt.Println("=== Original bin list ===")
	binList.PrintAll()

	// Create storage instance via DI
	fs := f.NewFS()
	var storage s.Store = s.New(fs, "data/bins.json")

	// Save bins to JSON file
	fmt.Println("\n=== Saving bins to JSON ===")
	if err := storage.SaveBins(binList); err != nil {
		log.Printf("Error saving bins: %v", err)
	} else {
		fmt.Printf("Bins saved to: %s\n", storage.GetStoragePath())
	}

	// Load bins from JSON file
	fmt.Println("\n=== Loading bins from JSON ===")
	loadedBinList, err := storage.LoadBins()
	if err != nil {
		log.Printf("Error loading bins: %v", err)
	} else {
		fmt.Println("Loaded bins from storage:")
		loadedBinList.PrintAll()
	}

	// Example of searching for a bin by ID
	if foundBin, exists := binList.GetByID("bin002"); exists {
		fmt.Printf("\nFound bin: %s (Private: %t)\n", foundBin.Name, foundBin.Private)
	}
}
