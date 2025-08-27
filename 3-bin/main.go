package main

import (
	"fmt"
	"time"
)

// Bin represents a structure for storing bin information
type Bin struct {
	id        string
	private   bool
	createdAt time.Time
	name      string
}

// BinList represents a list of Bin
type BinList struct {
	bins []Bin
}

// NewBin creates a new instance of Bin
func NewBin(id, name string, private bool) *Bin {
	return &Bin{
		id:        id,
		private:   private,
		createdAt: time.Now(),
		name:      name,
	}
}

// NewBinList creates a new instance of BinList
func NewBinList() *BinList {
	return &BinList{
		bins: make([]Bin, 0),
	}
}

// AddBin adds a bin to the list
func (bl *BinList) AddBin(bin Bin) {
	bl.bins = append(bl.bins, bin)
}

// GetBinByID returns a bin by ID
func (bl *BinList) GetBinByID(id string) (*Bin, bool) {
	for _, bin := range bl.bins {
		if bin.id == id {
			return &bin, true
		}
	}
	return nil, false
}

// PrintAllBins prints all bins to console
func (bl *BinList) PrintAllBins() {
	fmt.Println("List of all bins:")
	for i, bin := range bl.bins {
		fmt.Printf("%d. ID: %s, Name: %s, Private: %t, Created: %s\n",
			i+1, bin.id, bin.name, bin.private, bin.createdAt.Format("2006-01-02 15:04:05"))
	}
}

func main() {
	// Create a new list
	binList := NewBinList()

	// Create several example bins
	bin1 := NewBin("bin001", "My first bin", false)
	bin2 := NewBin("bin002", "Private bin", true)
	bin3 := NewBin("bin003", "Work bin", false)

	// Add bins to the list
	binList.AddBin(*bin1)
	binList.AddBin(*bin2)
	binList.AddBin(*bin3)

	// Print all bins
	binList.PrintAllBins()

	// Example of searching for a bin by ID
	if foundBin, exists := binList.GetBinByID("bin002"); exists {
		fmt.Printf("\nFound bin: %s (Private: %t)\n", foundBin.name, foundBin.private)
	}
}
