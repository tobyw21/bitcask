package storage

/*
	disk.go defines KV entry stores on the disk
*/

type KVEntry[T any] struct {
	Crc       string
	TimeStamp int64
	KeySz     uintptr
	ValueSz   uintptr
	Key       string
	Value     T
}
