package storage

/*
	catalog.go generates oid to KV set to save the files on the
		file systems, a file system can have multiple KV stores
		it simulates Postgres' Oid system, a catalog file will
		be saved on disk to keep a table of information about
		used Oids and its mapping to corresponding KV store
*/

type Oid uint32

const NextDataOid Oid = 1
const NextKVOid Oid = 1000

func KVOidExistsInCatalog(oid Oid) bool {

	return false
}

func GetNewKVOid() Oid {
	var new_kv_oid Oid
	return new_kv_oid
}

func GetNewDataOid() Oid {
	var new_data_oid Oid
	return new_data_oid
}
