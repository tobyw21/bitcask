package storage

import "time"

/*
	disk.go defines KV entry stores on the disk
*/

type KVEntry[T any] struct {
	Crc 		string
	TimeStamp 	time.Time
	KeySz	  	uint32
	ValueSz		uint32
	Key 		string
	Value 		T
}
