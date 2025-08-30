package bins

import (
	"fmt"
	"time"
)

// Bin represents a structure for storing bin information
type Bin struct {
	ID        string    `json:"id"`
	Private   bool      `json:"private"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
}

// BinList represents a list of Bin
type BinList struct {
	Bins []Bin `json:"bins"`
}

// NewBin creates a new instance of Bin
func NewBin(id, name string, private bool) *Bin {
	return &Bin{
		ID:        id,
		Private:   private,
		CreatedAt: time.Now(),
		Name:      name,
	}
}

// NewList creates a new instance of BinList
func NewList() *BinList {
	return &BinList{Bins: make([]Bin, 0)}
}

// Add adds a bin to the list
func (bl *BinList) Add(bin Bin) {
	bl.Bins = append(bl.Bins, bin)
}

// GetByID returns a bin by ID
func (bl *BinList) GetByID(id string) (*Bin, bool) {
	for i := range bl.Bins {
		if bl.Bins[i].ID == id {
			return &bl.Bins[i], true
		}
	}
	return nil, false
}

// PrintAll prints all bins to console
func (bl *BinList) PrintAll() {
	fmt.Println("List of all bins:")
	for i, bin := range bl.Bins {
		fmt.Printf("%d. ID: %s, Name: %s, Private: %t, Created: %s\n",
			i+1, bin.ID, bin.Name, bin.Private, bin.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}
