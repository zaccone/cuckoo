package cuckoo

import (
	"errors"
	"fmt"
)

type hashFn func(key uint64, size uint64) uint64

// ErrTableEmpty is a dedicated error indicating the element cannot be added due to collisions
var ErrTableEmpty = errors.New("Cuckoo table is full")

func yinHash(key uint64, size uint64) uint64 {
	return key % size
}

func yangHash(key uint64, size uint64) uint64 {
	return (key / size) % size
}

// Item is an object stores in the Cuckoo hash object
type Item struct {
	Key   uint64
	Value string
}

func (i *Item) String() string {
	return fmt.Sprintf("Item key: %v, Value: %v", i.Key, i.Value)
}

// Cuckoo is a structure representing whole hash object
// For now lets make it of key integer and value string
type Cuckoo struct {
	size uint64
	yin  []*Item
	yang []*Item

	yinHash  hashFn
	yangHash hashFn
}

// New returns new Cuckoo hash object with the size initialized via the parameter
func New(size uint64, yinFn hashFn, yangFn hashFn) *Cuckoo {
	return &Cuckoo{size, make([]*Item, size), make([]*Item, size), yinFn, yangFn}
}

// NewDefaultHash returns a Cuckoo object with internal tables set to size and default hash functions: yinHash and yangHash
func NewDefaultHash(size uint64) *Cuckoo {
	return &Cuckoo{size, make([]*Item, size), make([]*Item, size), yinHash, yangHash}
}

// Insert inserts or replaces element stored under the key.
func (c *Cuckoo) Insert(key uint64, value string) bool {

	// try first of two possible locations
	yinH := c.yinHash(key, c.size)
	if c.yin[yinH] == nil {
		c.yin[yinH] = &Item{Key: key, Value: value}
		return true
	}

	// try second of two possible locations
	yangH := c.yangHash(key, c.size)
	if c.yang[yangH] == nil {
		c.yang[yangH] = &Item{Key: key, Value: value}
		return true
	}

	// try misplacing element from yin and see if that is possible. If so, enter the element into yin

	if c.Insert(c.yin[yinH].Key, c.yin[yinH].Value) {
		c.yin[yinH] = &Item{Key: key, Value: value}
		return true
	}
	// try misplacing element from yang and see if that is possible. If so, enter the element into yang

	if c.Insert(c.yang[yangH].Key, c.yang[yangH].Value) {
		c.yin[yangH] = &Item{Key: key, Value: value}
		return true
	}

	// Apparently we couldn't find enough room in either of the hashes, return nil, as the table is considered 'full'
	return false // for now, lol

}

// Lookup finds element stored as key and returns pointer to this element or nil if the element does not exist.
func (c *Cuckoo) Lookup(key uint64) *Item {
	yin := c.yin[c.yinHash(key, c.size)]
	if ok(yin, key) {
		return yin
	}

	yang := c.yang[yangHash(key, c.size)]
	if ok(yang, key) {
		return yang
	}
	return nil
}

// Delete removes element from the hashtable and returns true if the element was stored, false otherwise.
func (c *Cuckoo) Delete(key uint64) bool {
	return false
}

// Debug prints values from both yin ang yang tables
func (c *Cuckoo) Debug() {
	fmt.Println("Yin: ")
	for idx, item := range c.yin {
		str := "nil"
		if item != nil {
			str = item.String()
		}
		fmt.Printf("%d: %s", idx, str)
	}
	fmt.Println("")

	fmt.Println("Yang: ")
	for idx, item := range c.yang {
		str := "nil"
		if item != nil {
			str = item.String()
		}
		fmt.Printf("%d: %s", idx, str)
	}
	fmt.Println("")
}

// helper function to extract Item if it really matches
func ok(item *Item, key uint64) bool {
	return item != nil && item.Key == key
}
