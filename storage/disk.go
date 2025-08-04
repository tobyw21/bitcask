package storage

/*
	disk.go defines KV entry stores on the disk
*/

type Header struct {
	HeaderSize uint64
}

type KVEntry struct {
	Crc       string
	TimeStamp int64
	KeySz     int64
	ValueSz   int64
	Key       string
	Value     interface{}
}
