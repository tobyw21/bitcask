package fs

/*
	catalog.go generates oid to KV set to save the files on the
		file systems, a file system can have multiple KV stores
		it simulates Postgres' Oid system, a catalog file will
		be saved on disk to keep a table of information about
		used Oids and its mapping to corresponding KV store
*/

type Oid uint32

func OidExistsInCatalog(oid Oid) bool {
	return true
}

func GetNewOid() Oid {
	return 42
}
