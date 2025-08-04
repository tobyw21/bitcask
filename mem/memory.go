package mem

/*
	memory.go defines data structure in the memory
*/

import (
	"github.com/tobyw21/bitcask/include"
	"github.com/tobyw21/bitcask/vfd"
)

type KeyDir struct {
	FileId    include.Oid
	ValueSz   int64
	ValuePos  int64
	TimeStamp int64
}

func NewKeyDir(fileid include.Oid, valsize int64, valpos int64, timestmp int64) KeyDir {
	return KeyDir{
		FileId:    fileid,
		ValueSz:   valsize,
		ValuePos:  valpos,
		TimeStamp: timestmp,
	}
}

func WriteHint(vfdmgr vfd.VfdManager,  kdm map[string]KeyDir) {

}

func ReadHint() map[string]KeyDir {
	return nil
}
