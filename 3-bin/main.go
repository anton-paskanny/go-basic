package main

import (
	b "demo/bin/bins"
	"fmt"
)

func main() {
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
	binList.PrintAll()

	// Example of searching for a bin by ID
	if foundBin, exists := binList.GetByID("bin002"); exists {
		fmt.Printf("\nFound bin: %s (Private: %t)\n", foundBin.Name(), foundBin.Private())
	}
}
