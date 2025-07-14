package bitcask

import (
	"encoding/gob"
	"fmt"
	"reflect"
	"time"

	"github.com/tobyw21/bitcask/include"
	"github.com/tobyw21/bitcask/mem"
	"github.com/tobyw21/bitcask/storage"
	"github.com/tobyw21/bitcask/vfd"
)

/*
	bitstore.go is the API exported to user to do the get, set, remove operations

*/

type BitStoreManager[T any] struct {
	kvoid     include.Oid
	name      string
	currFile  include.Oid
	nextFree  int64
	catMgr    *storage.Catalog
	vfdMgr    *vfd.VfdManager
	keydirMap map[string]mem.KeyDir
}

func BitStore[T any](name string) *BitStoreManager[T] {

	v := vfd.NewVfdMgr()
	c, err := storage.CatalogRead(v, "data/catalog")

	if err != nil {
		panic("can't read catalog file!")
	}
	kvoid := c.SetKvStoreMap(name)
	curroid := c.GetDatNextOid(kvoid)
	c.AppendDataStoreMap(kvoid, curroid)

	// kd := mem.ReadHint()
	kd := make(map[string]mem.KeyDir)

	return &BitStoreManager[T]{
		kvoid:     kvoid,
		name:      name,
		currFile:  curroid,
		catMgr:    c,
		vfdMgr:    v,
		keydirMap: kd,
	}
}

func (b *BitStoreManager[T]) Get(key string) T {
	var t T
	return t
}

func (b *BitStoreManager[T]) Set(key string, value T) error {

	filepath := fmt.Sprintf("data/%d/%d", b.kvoid, b.currFile)
	vid, err := b.vfdMgr.VfdOpen(filepath)
	defer b.vfdMgr.VfdClose(vid)

	if err != nil {
		return err
	}

	w := &vfd.VfdWriter{Vfdid: vid, Offset: b.nextFree, Vfdmgr: b.vfdMgr}
	var t int64 = time.Now().Unix()
	vs := reflect.TypeOf(value).Size()
	kd := mem.NewKeyDir(b.currFile, vs, b.nextFree, t)
	b.keydirMap[key] = kd

	kve := storage.KVEntry[T]{
		Crc:       "something",
		TimeStamp: t,
		KeySz:     reflect.TypeOf(key).Size(),
		ValueSz:   vs,
		Key:       key,
		Value:     value,
	}

	enc := gob.NewEncoder(w)
	err = enc.Encode(&kve)
	return err
}
