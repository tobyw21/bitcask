package storage

import (
	"os"

	vfd "github.com/tobyw21/bitcask/vfd"
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
	KvStoreMap  map[string]Oid // kv name : Oid
}

func Serialise(c Catalog) []byte {
	return make([]byte, 0)
}

func Deserialise(b []byte) Catalog {
	return Catalog{}
}

func CatalogLoad(vfdmgr *vfd.VfdManager, path string) *Catalog {

	catalog := Catalog{
		NextDataOid: 0,
		NextKvOid:   0,
		KvStoreMap:  make(map[string]Oid),
	}

	if _, err := os.Stat(path); err != nil {

		return &catalog
	} else {
		vid, err := vfdmgr.VfdOpen(path)

		if err != nil {
			panic("Unable to open file " + path)
		}
		// load catalog from file

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

func (c *Catalog) CatalogWrite(vfdmgr *vfd.VfdManager, path string) {

	vid, err := vfdmgr.VfdOpen(path)
	if err != nil {
		panic("Unable to open file " + path)
	}

	vfdmgr.VfdClose(vid)
}
