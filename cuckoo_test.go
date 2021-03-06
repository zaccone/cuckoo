package cuckoo

import "testing"

const TableSize = 11

func itemCompare(i1 *Item, i2 *Item) bool {
	if i1 == nil && i2 == nil {
		return true
	}

	return i1.Key == i2.Key && i1.Value == i2.Value
}

func TestLookupNil(t *testing.T) {
	cuckoo := NewDefaultHash(TableSize)

	for k := uint64(0); k < 100; k++ {
		if cuckoo.Lookup(k) != nil {
			t.Errorf("Expected value nil for key: %d", k)
		}
	}

}

func TestInsertOnEmpty(t *testing.T) {
	cuckoo := NewDefaultHash(TableSize)

	item := &Item{Key: 13, Value: "marek"}

	if cuckoo.Lookup(item.Key) != nil {
		t.Errorf("Expected nil during first Cuckoo.Lookup()")
	}

	if cuckoo.Insert(item.Key, item.Value) == false {
		t.Errorf("Expected true while inserting, got false")
	}

	stored := cuckoo.Lookup(item.Key)

	if stored == nil {
		t.Errorf("Expected !nil during second Cuckoo.Lookup()")
	}

	if stored.Key != item.Key || stored.Value != item.Value {
		t.Errorf("Item mismatch item: %s vs stored: %s", item.String(), stored.String())
	}

}

func TestInsertWithSameKey(t *testing.T) {

	cuckoo := NewDefaultHash(TableSize)

	i1 := &Item{Key: 3, Value: "three"}
	i2 := &Item{Key: 36, Value: "thirty six"}

	if cuckoo.yinHash(i1.Key, TableSize) != cuckoo.yangHash(i2.Key, TableSize) {
		t.Errorf("Key equality is a prerequisite for collision resolution test: yinHash %d, yangHash: %d",
			cuckoo.yinHash(i1.Key, TableSize), cuckoo.yangHash(i2.Key, TableSize))
	}

	cuckoo.Insert(i1.Key, i1.Value)
	cuckoo.Insert(i2.Key, i2.Value)

	ci1 := cuckoo.Lookup(i1.Key)

	if itemCompare(i1, ci1) == false {
		t.Errorf("Item 1 mismatch: original item %s, fetched item: %s", i1.String(), ci1.String())
	}

	ci2 := cuckoo.Lookup(i2.Key)

	if itemCompare(i2, ci2) == false {
		t.Errorf("Item 2 mismatch: original item %s, fetched item: %s", i2.String(), ci2.String())
	}

}

/* TestInsertWithEviction tests inserting and lookups where the eviction must happen. For the standard hash function we will have following values hashing to values:

Key  YinFn YangFn
105	 6	   9
3	 3	   0
36	 3	   3
39	 6	   3

*/
func TestInsertWithEviction(t *testing.T) {
	cuckoo := NewDefaultHash(TableSize)

	/*
		We trigger this order of inserts to trigger cuckko eviction. At first Item with key 105 will be hashed to yin[6], however, later
		after adding 39, it will be rehashed to yang[9].


	*/
	i1 := &Item{Key: 105, Value: "hundred and five"}
	i2 := &Item{Key: 3, Value: "three"}
	i3 := &Item{Key: 36, Value: "thirty six"}
	i4 := &Item{Key: 39, Value: "thirty nine"}

	cuckoo.Insert(i1.Key, i1.Value)
	cuckoo.Insert(i2.Key, i2.Value)
	cuckoo.Insert(i3.Key, i3.Value)
	cuckoo.Insert(i4.Key, i4.Value)

	ci1 := cuckoo.Lookup(i1.Key)

	if itemCompare(i1, ci1) == false {
		t.Errorf("Item 1 mismatch: original item %s, fetched item: %s", i1.String(), ci1.String())
	}

	ci2 := cuckoo.Lookup(i2.Key)

	if itemCompare(i2, ci2) == false {
		t.Errorf("Item 2 mismatch: original item %s, fetched item: %s", i2.String(), ci2.String())
	}

	ci3 := cuckoo.Lookup(i3.Key)

	if itemCompare(i3, ci3) == false {
		t.Errorf("Item 3 mismatch: original item %s, fetched item: %s", i3.String(), ci3.String())
	}

	ci4 := cuckoo.Lookup(i4.Key)

	if itemCompare(i4, ci4) == false {
		t.Errorf("Item 4 mismatch: original item %s, fetched item: %s", i4.String(), ci4.String())
	}

	/* 	Ensure that items are in the correct  indexes in the internal tables */

	var it *Item

	it = cuckoo.yin[3]

	if it == nil {
		t.Errorf("Expected to have non nil *Item at cucko.yin[3]")
	}

	if itemCompare(it, i2) == false {
		t.Errorf("Item mismatch, expected got: %s, expected: %s", it.String(), i2.String())
	}

	it = cuckoo.yin[6]
	if it == nil {
		t.Errorf("Expected to have non nil *Item at cucko.yin[6]")
	}

	if itemCompare(it, i4) == false {
		t.Errorf("Item mismatch, expected got: %s, expected: %s", it.String(), i2.String())
	}

	it = cuckoo.yang[3]
	if it == nil {
		t.Errorf("Expected to have non nil *Item at cucko.yang[3]")
	}

	if itemCompare(it, i3) == false {
		t.Errorf("Item mismatch, expected got: %s, expected: %s", it.String(), i2.String())
	}

	it = cuckoo.yang[9]

	if it == nil {
		t.Errorf("Expected to have non nil *Item at cucko.yang[9]")
	}

	if itemCompare(it, i1) == false {
		t.Errorf("Item mismatch, expected got: %s, expected: %s", it.String(), i2.String())
	}

	// Usefull to debug when you run
	// $ go test -v
	t.Log("Inspecting Cuckoo.yin")
	for idx, item := range cuckoo.yin {
		if item != nil {
			t.Logf("At index %d found %s", idx, item.String())
		}
	}

	t.Log("Inspecting Cuckoo.yang")
	for idx, item := range cuckoo.yang {
		if item != nil {
			t.Logf("At index %d found %s", idx, item.String())
		}
	}

}

func TestEmptyLookup(t *testing.T) {
	cuckoo := NewDefaultHash(TableSize)

	// Value 3 hashes to
	i1 := &Item{Key: 3, Value: "three"}
	i2 := &Item{Key: 36, Value: "thirty six"}
	i3 := &Item{Key: 39, Value: "thirty nine"}
	i4 := &Item{Key: 105, Value: "hundred and five"}

	cuckoo.Insert(i1.Key, i1.Value)
	cuckoo.Insert(i2.Key, i2.Value)
	cuckoo.Insert(i3.Key, i3.Value)
	cuckoo.Insert(i4.Key, i4.Value)

	if cuckoo.Lookup(1) != nil {
		t.Error("Expected to not found any element.")
	}

}
