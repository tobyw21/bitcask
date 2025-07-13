package storage

import (
	"bytes"
	"container/list"
	"encoding/gob"

	"github.com/tobyw21/bitcask/include"
	vfd "github.com/tobyw21/bitcask/vfd"
)

/*
	catalog.go generates oid to KV set to save the files on the
		file systems, a file system can have multiple KV stores
		it simulates Postgres' Oid system, a catalog file will
		be saved on disk to keep a table of information about
		used Oids and its mapping to corresponding KV store
*/

type Catalog struct {
	NextDataOid  include.Oid
	NextKvOid    include.Oid
	KvStoreMap   map[string]include.Oid    // kv name : Oid
	DataStoreMap map[include.Oid]*list.List // kv oid: list of files
}

func (c *Catalog) GetKvNextOid() include.Oid {
	c.NextKvOid++
	return c.NextKvOid
}

func (c *Catalog) GetDatNextOid() include.Oid {
	c.NextDataOid++
	return c.NextDataOid
}

func CatalogLoad(vfdmgr *vfd.VfdManager, path string) *Catalog {

}

func (c *Catalog) CatalogWrite(vfdmgr *vfd.VfdManager, path string) error {

	vid, err := vfdmgr.VfdOpen(path)
	if err != nil {
		return err
	}
	defer vfdmgr.VfdClose(vid)
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	enc.Encode(c)

	_, err = vfdmgr.VfdWrite(vid, buf.Bytes(), 0)

	return err
}
