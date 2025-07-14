package mem

/*
	memory.go defines data structure in the memory
*/

import (
	"github.com/tobyw21/bitcask/include"
)

type KeyDir struct {
	FileId    include.Oid
	ValueSz   uintptr
	ValuePos  int64
	TimeStamp int64
}

func NewKeyDir(fileid include.Oid, valsize uintptr, valpos int64, timestmp int64) KeyDir {
	return KeyDir{
		FileId:    fileid,
		ValueSz:   valsize,
		ValuePos:  valpos,
		TimeStamp: timestmp,
	}
}

func WriteHint() {

}

func ReadHint() map[string]KeyDir {
	return nil
}
