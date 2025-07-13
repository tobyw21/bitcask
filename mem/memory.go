package mem

/*
	memory.go defines data structure in the memory
*/

import (
	"time"

	"github.com/tobyw21/bitcask/include"
)

type TimeStampType time.Time

type KeyDir struct {
	FileId    include.Oid
	ValueSz   uint32
	ValuePos  uint32
	TimeStamp TimeStampType
}

func NewKeyDir(fileid include.Oid, valsize uint32, valpos uint32, timestmp TimeStampType) KeyDir {
	return KeyDir{
		FileId:    fileid,
		ValueSz:   valsize,
		ValuePos:  valpos,
		TimeStamp: timestmp,
	}
}
