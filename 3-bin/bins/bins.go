package bins

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

// NewList creates a new instance of BinList
func NewList() *BinList {
	return &BinList{bins: make([]Bin, 0)}
}

// Add adds a bin to the list
func (bl *BinList) Add(bin Bin) {
	bl.bins = append(bl.bins, bin)
}

// GetByID returns a bin by ID
func (bl *BinList) GetByID(id string) (*Bin, bool) {
	for i := range bl.bins {
		if bl.bins[i].id == id {
			return &bl.bins[i], true
		}
	}
	return nil, false
}

// PrintAll prints all bins to console
func (bl *BinList) PrintAll() {
	fmt.Println("List of all bins:")
	for i, bin := range bl.bins {
		fmt.Printf("%d. ID: %s, Name: %s, Private: %t, Created: %s\n",
			i+1, bin.id, bin.name, bin.private, bin.createdAt.Format("2006-01-02 15:04:05"))
	}
}

// ID returns bin id
func (b *Bin) ID() string { return b.id }

// Name returns bin name
func (b *Bin) Name() string { return b.name }

// Private returns whether bin is private
func (b *Bin) Private() bool { return b.private }

// CreatedAt returns creation time
func (b *Bin) CreatedAt() time.Time { return b.createdAt }
