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

type BitStoreManager struct {
	kvoid     include.Oid
	name      string
	currFile  include.Oid
	nextFree  int64
	catMgr    *storage.Catalog
	vfdMgr    *vfd.VfdManager
	keydirMap map[string]mem.KeyDir
}

func BitStore(name string) *BitStoreManager {

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

	return &BitStoreManager{
		kvoid:     kvoid,
		name:      name,
		currFile:  curroid,
		catMgr:    c,
		vfdMgr:    v,
		keydirMap: kd,
	}
}

func (b *BitStoreManager) Get(key string) (interface{}, error) {
	return nil, nil
}

func (b *BitStoreManager) Set(key string, value interface{}) error {

	filepath := fmt.Sprintf("data/%d/%d", b.kvoid, b.currFile)
	vid, err := b.vfdMgr.VfdOpen(filepath)
	// defer b.vfdMgr.VfdClose(vid)

	if err != nil {
		return err
	}

	w := &vfd.VfdWriter{Vfdid: vid, Offset: b.nextFree, Vfdmgr: b.vfdMgr}
	

	var t int64 = time.Now().Unix()
	vs := int64(reflect.TypeOf(value).Size())
	kd := mem.NewKeyDir(b.currFile, vs, b.nextFree, t)
	b.keydirMap[key] = kd

	kve := storage.KVEntry{
		Crc:       "something",
		TimeStamp: t,
		KeySz:     int64(len(key)), // int64(reflect.TypeOf(key).Size()),
		ValueSz:   vs,
		Key:       key,
		Value:     value,
	}

	b.nextFree += int64(reflect.TypeOf(kve).Size())

	
	return err
}

func (b *BitStoreManager) Delete(key string, value interface{}) error {
	return nil
}

func (b *BitStoreManager) Close() error {
	err := b.catMgr.CatalogWrite(b.vfdMgr, "data/catalog")

	if err != nil {
		return err
	}

	mem.WriteHint(*b.vfdMgr, b.keydirMap)

	b.vfdMgr.VfdClean()
	return nil
}
