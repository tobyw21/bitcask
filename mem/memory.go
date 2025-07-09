package mem

/*
	memory.go defines data structure in the memory
*/

import (
	"time"
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
