package bitcask

/*
	bitstore.go is the API exported to user to do the get, set, remove operations

*/

import (
	"fmt"

	mem "github.com/tobyw21/bitcask/mem"
	st "github.com/tobyw21/bitcask/storage"
	vfd "github.com/tobyw21/bitcask/vfd"
)

type BitStore[T any] struct {
	KvName         string
	vfdManager     *vfd.VfdManager
	catalogManager *st.Catalog
	keyDirManager  map[string]mem.KeyDir
}

func NewBitStore[T any](name string) *BitStore[T] {

	// lets just hard code path for now
	path := "data/catalog"

	vfdmgr := vfd.NewVfdMgr()
	catalog := st.CatalogLoad(vfdmgr, path)
	kd := make(map[string]mem.KeyDir)

	return &BitStore[T]{KvName: name, vfdManager: vfdmgr, catalogManager: catalog, keyDirManager: kd}
}

func (b *BitStore[T]) Get(key string) T {

	dir := b.keyDirManager[key]

	// use fildId, value position to find data
	path := fmt.Sprintf("data/%d/%d", b.catalogManager.KvStoreMap[b.KvName], dir.FileId)

	vid, err := b.vfdManager.VfdOpen(path)
	defer b.vfdManager.VfdClose(vid)
	if err != nil {
		panic("")
	}
	

	buf := make([]byte, b.keyDirManager[key].ValueSz)
	b.vfdManager.VfdRead(vid, buf, int64(b.keyDirManager[key].ValuePos))

	// decode buf to KVEntry...

	kve := st.KVEntry[T]{}
	// check crc update other meta data...

	// return data
	var t T = kve.Value
	return t
}
