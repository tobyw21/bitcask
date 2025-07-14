package bitcask

import (
	"testing"
)

func TestBitStore(t *testing.T) {
	b := BitStore[int]("store1")

	err := b.Set("apple", 1)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

}
