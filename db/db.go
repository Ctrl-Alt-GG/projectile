package db

import (
	"github.com/Ctrl-Alt-GG/projectile/model"
	"slices"
	"sync"
)

var ( // this is a lame, locking array
	db     []model.GameServer
	dbLock sync.RWMutex
)

func StoreUpdate(update model.GameServer) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	idx := slices.IndexFunc(db, func(server model.GameServer) bool {
		return server.Address == update.Address
	})

	if idx == -1 {
		db = append(db, update)
	} else {
		db[idx] = update
	}

	return nil
}

func DeleteServer(address string) (int, error) {
	dbLock.Lock()
	defer dbLock.Unlock()

	lenBefore := len(db)

	db = slices.DeleteFunc(db, func(server model.GameServer) bool {
		return server.Address == address
	})

	return lenBefore - len(db), nil
}

func GetAll() []model.GameServer {
	dbLock.RLock()
	defer dbLock.RUnlock()

	dbCopy := make([]model.GameServer, len(db))

	for i, val := range db {
		dbCopy[i] = val.Copy()
	}

	return dbCopy
}
