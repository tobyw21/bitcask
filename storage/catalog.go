package storage

import (
	"bytes"
	"encoding/gob"
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

func CatalogLoad(vfdmgr *vfd.VfdManager, path string) *Catalog {

	catalog := Catalog{
		NextDataOid: 0,
		NextKvOid:   0,
		KvStoreMap:  make(map[string]Oid),
	}
	fi, err := os.Stat(path)

	if err == nil {
		vid, err := vfdmgr.VfdOpen(path)

		defer vfdmgr.VfdClose(vid)
		if err != nil {
			panic("Unable to open file " + path)
		}
		// load catalog from file

		s := fi.Size()
		tmp_buf := make([]byte, s)
		vfdmgr.VfdRead(vid, tmp_buf, 0)

		var buf *bytes.Buffer = bytes.NewBuffer(tmp_buf)
		dec := gob.NewDecoder(buf)

		dec.Decode(&catalog)

		return &catalog
	}

	return &catalog

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
	defer vfdmgr.VfdClose(vid)
	if err != nil {
		panic("Unable to open file " + path)
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	enc.Encode(c)

	vfdmgr.VfdWrite(vid, buf.Bytes(), 0)

}
