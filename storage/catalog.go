package storage

import (
	"encoding/gob"
	"os"

	"github.com/tobyw21/bitcask/include"
	"github.com/tobyw21/bitcask/vfd"
)

/*
	catalog.go generates oid to KV set to save the files on the
		file systems, a file system can have multiple KV stores
		it simulates Postgres' Oid system, a catalog file will
		be saved on disk to keep a table of information about
		used Oids and its mapping to corresponding KV store
*/

type Catalog struct {
	kvOid        include.Oid
	kvStoreMap   map[string]include.Oid        // kv name : Oid
	dataStoreMap map[include.Oid][]include.Oid // kv oid: list of files
}

func (c *Catalog) GetKvNextOid() include.Oid {
	c.kvOid++
	return c.kvOid
}

func (c *Catalog) GetDatNextOid(kvoid include.Oid) include.Oid {

	if l := len(c.dataStoreMap[kvoid]) - 1; l > 0 {
		return c.dataStoreMap[kvoid][l] + 1
	}
	// if no data store file
	// oid starts from 1 for each store
	return 1
}

func (c *Catalog) AppendDataStoreMap(kvoid include.Oid, dataoid include.Oid) {
	c.dataStoreMap[kvoid] = append(c.dataStoreMap[kvoid], dataoid)
}

func (c *Catalog) SetKvStoreMap(name string) include.Oid {

	if _, found := c.kvStoreMap[name]; !found {
		oid := c.GetKvNextOid()
		c.kvStoreMap[name] = oid
		return oid
	}

	return c.kvStoreMap[name]
}

func CatalogRead(vfdmgr *vfd.VfdManager, path string) (*Catalog, error) {
	var c Catalog = Catalog{
		kvOid:        0,
		kvStoreMap:   make(map[string]include.Oid),
		dataStoreMap: make(map[include.Oid][]include.Oid),
	}

	if f, err := os.Stat(path); err == nil {
		// if file exists
		// unmarshal
		_ = f
		vid, err := vfdmgr.VfdOpen(path)
		if err != nil {
			return nil, err
		}
		defer vfdmgr.VfdClose(vid)
		r := &vfd.VfdReader{Vfdid: vid, Offset: 0, Vfdmgr: vfdmgr}
		enc := gob.NewDecoder(r)
		enc.Decode(&c)

	}

	return &c, nil
}

func (c *Catalog) CatalogWrite(vfdmgr *vfd.VfdManager, path string) error {

	vid, err := vfdmgr.VfdOpen(path)
	if err != nil {
		return err
	}
	defer vfdmgr.VfdClose(vid)

	w := &vfd.VfdWriter{Vfdid: vid, Offset: 0, Vfdmgr: vfdmgr}
	enc := gob.NewEncoder(w)

	err = enc.Encode(c)

	return err
}
