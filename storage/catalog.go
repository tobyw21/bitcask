package storage

import (
	"os"
)

/*
	catalog.go generates oid to KV set to save the files on the
		file systems, a file system can have multiple KV stores
		it simulates Postgres' Oid system, a catalog file will
		be saved on disk to keep a table of information about
		used Oids and its mapping to corresponding KV store
*/

type Oid uint32

type Catalog struct {
	NextDataOid Oid
	NextKvOid   Oid
	OidList     []Oid
	KvStoreMap  []map[string]Oid
}

func LoadCatalog(vfdmgr *VfdManager) *Catalog {
	var path string = "data/catalog"
	if _, err := os.Stat(path); err != nil {

		return &Catalog{
			NextDataOid: 0,
			NextKvOid:   0,
			KvStoreMap:  make([]map[string]Oid, 0),
		}
	} else {
		vid, err := vfdmgr.VfdOpen(path)

		if err != nil {
			panic("Unable to open file " + path)
		}
		var catalog Catalog
		// := Catalog{
		// 	NextDataOid: 0,
		// 	NextKvOid:   0,
		// 	KvStoreMap:  make([]map[string]Oid, 0),
		// }

		// buffer := make([]byte, 56)
		// vfdmgr.VfdRead(vid, &buffer, 0)

		vfdmgr.VfdClose(vid)

		return &catalog
	}
}

func (c *Catalog) GetKvNextOid() Oid {
	c.NextKvOid++
	return c.NextKvOid
}

func (c *Catalog) GetDatNextOid() Oid {
	c.NextDataOid++
	return c.NextDataOid
}

func (c *Catalog) CatalogWrite(vfdmgr *VfdManager) {
	var path string = "data/catalog"
	vid, err := vfdmgr.VfdOpen(path)
	if err != nil {
		panic("Unable to open file " + path)
	}

	// vfdmgr.VfdWrite(vid, buffer, 0)

	vfdmgr.VfdClose(vid)
}
