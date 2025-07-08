package mem

/*
	memory.go defines data structure in the memory
*/

import (
	"time"

	st "github.com/tobyw21/bitcask/storage"
	vfd "github.com/tobyw21/bitcask/vfd"
)

type Oid uint32
type TimeStampType time.Time

type KeyDir struct {
	FileId    Oid
	ValueSz   uint32
	ValuePos  uint32
	TimeStamp TimeStampType
}

func NewKeyDir(fileid Oid, valsize uint32, valpos uint32, timestmp TimeStampType) KeyDir {
	return KeyDir{
		FileId:    fileid,
		ValueSz:   valsize,
		ValuePos:  valpos,
		TimeStamp: timestmp,
	}
}

func HintWrite(vfdmgr *vfd.VfdManager, c *st.Catalog, kvname string, kd []KeyDir) {
	// var oid Oid = Oid(c.KvStoreMap[kvname])
	// var path string = fmt.Sprintf("data/%d/%d.hint", oid, oid)

}
