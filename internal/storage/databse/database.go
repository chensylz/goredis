package databse

import (
	"sync"

	"github.com/chensylz/goredis/internal/storage"
	"github.com/chensylz/goredis/internal/storage/memory"
)

type Database struct {
	dbs []storage.DB
	sync.Mutex
}

func (d *Database) Get(index uint8) storage.DB {
	return d.dbs[index]
}

func New(size int) *Database {
	db := &Database{
		dbs: make([]storage.DB, size),
	}
	for i := 0; i < size; i++ {
		db.dbs[i] = memory.NewSyncDict(i)
	}
	return db
}
