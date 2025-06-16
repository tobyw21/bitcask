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
