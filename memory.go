package bitcask

import (
	"time"
)

type KeyDir struct {
	FileId    uint
	ValueSz   uint
	ValuePos  uint
	TimeStamp time.Time
}

func NewKeyDir() KeyDir {
	panic("TODO")
}

func (kd *KeyDir) Update(args ...[]byte) {
	panic("TODO Pleace holder")
}

func (kd *KeyDir) Remove(args ...[]byte) {
	panic("TODO Pleace holder")
}
