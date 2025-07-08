package storage

import (
	"testing"

	vfd "github.com/tobyw21/bitcask/vfd"
)

func TestCreateCatalog(t *testing.T) {
	vfdmgr := vfd.NewVfdMgr()
	var path string = "../examples/catalog"

	c := CatalogLoad(vfdmgr, path)
	// c.KvStoreMap["test_store1"] = c.GetKvNextOid()
	// c.KvStoreMap["test_store2"] = c.GetKvNextOid()
	// c.KvStoreMap["test_store3"] = c.GetKvNextOid()

	// c.CatalogWrite(vfdmgr, path)

	t.Logf("%v", c)
}
