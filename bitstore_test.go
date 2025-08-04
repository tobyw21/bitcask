package bitcask

import (
	"testing"
)

func TestBitStore(t *testing.T) {
	b := BitStore("store1")

	err := b.Set("apple", 1)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	v, err := b.Get("apple")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Logf("returned value = %d", v)

	b.Close()

}
