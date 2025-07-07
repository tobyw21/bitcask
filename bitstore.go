package bitcask

/*
	bitstore.go is the API exported to user to do the get, set, remove operations

*/

type BitStore[T any] struct{}

func NewBitStore[T any](name string) *BitStore[T] {
	return &BitStore[T]{}
}

func (b *BitStore[T]) Get(key string) T {
	var t T
	return t
}