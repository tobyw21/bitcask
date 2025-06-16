package mem

/*
	memory.go defines KeyDir structure in the memory
*/

import (
	"time"
)

type Oid uint32
type KeyDir struct {
	FileId    Oid
	ValueSz   uint32
	ValuePos  uint32
	TimeStamp time.Time
}

func NewKeyDir(file_id Oid, value_size uint32, value_position uint32, time_stamp time.Time) KeyDir {
	return KeyDir{
		FileId:    file_id,
		ValueSz:   value_size,
		ValuePos:  value_position,
		TimeStamp: time_stamp,
	}
}
