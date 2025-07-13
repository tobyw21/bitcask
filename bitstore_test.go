package bitcask

import (
	"testing"
)

func TestBitStore(t *testing.T) {
	b := BitStore[int32]("bitstore1")

	b.Set("data1", 1)

	d := b.Get("data1")

	t.Logf("data get = %d", d)
}
