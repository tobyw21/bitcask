package bitcask

/*
	bitstore.go is the API exported to user to do the get, set, remove operations

*/

type BitStoreManager[T any] struct {
}

func BitStore[T any](name string) *BitStoreManager[T] {

}

func (b *BitStoreManager[T]) Get(key string) T {

}

func (b *BitStoreManager[T]) Set(key string, value T) error {

}
