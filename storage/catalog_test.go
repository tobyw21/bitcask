package storage

import (
	"testing"

	"github.com/tobyw21/bitcask/vfd"
)

func TestWriteCatalog(t *testing.T) {
	vfdmgr := vfd.NewVfdMgr()
	var path string = "../examples/catalog"

	c, err := CatalogRead(vfdmgr, path)

	if err != nil {
		t.FailNow()
		t.Log(err)
	}

	c.kvStoreMap["test_store1"] = c.GetKvNextOid()
	c.kvStoreMap["test_store2"] = c.GetKvNextOid()
	c.kvStoreMap["test_store3"] = c.GetKvNextOid()

	c.CatalogWrite(vfdmgr, path)
}

func TestReadCatalog(t *testing.T) {
	vfdmgr := vfd.NewVfdMgr()
	var path string = "../examples/catalog"

	c, err := CatalogRead(vfdmgr, path)
	if err != nil {
		t.FailNow()
		t.Log(err)
	}
	t.Logf("%v", c)
}
