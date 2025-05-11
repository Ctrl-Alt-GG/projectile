package db

import (
	"github.com/Ctrl-Alt-GG/projectile/model"
	"slices"
	"sync"
	"time"
)

var ( // this is a lame, locking array
	gameServers     []model.GameServer
	gameServersLock sync.RWMutex
)

func StoreUpdate(update model.GameServer) error {
	gameServersLock.Lock()
	defer gameServersLock.Unlock()

	idx := slices.IndexFunc(gameServers, func(server model.GameServer) bool {
		return server.Address == update.Address
	})

	if idx == -1 {
		gameServers = append(gameServers, update)
	} else {
		gameServers[idx] = update
	}

	return nil
}

func DeleteServer(address string) (int, error) {
	gameServersLock.Lock()
	defer gameServersLock.Unlock()

	lenBefore := len(gameServers)

	gameServers = slices.DeleteFunc(gameServers, func(server model.GameServer) bool {
		return server.Address == address
	})

	return lenBefore - len(gameServers), nil
}

func CleanupJob() (int, error) {
	gameServersLock.Lock()
	defer gameServersLock.Unlock()

	now := time.Now()

	lenBefore := len(gameServers)

	gameServers = slices.DeleteFunc(gameServers, func(server model.GameServer) bool {
		// delete entries older than a minute
		return now.Sub(server.LastUpdate) > time.Minute
	})

	return lenBefore - len(gameServers), nil
}

func GetAll() []model.GameServer {
	gameServersLock.RLock()
	defer gameServersLock.RUnlock()

	dbCopy := make([]model.GameServer, len(gameServers))

	for i, val := range gameServers {
		dbCopy[i] = val.Copy()
	}

	return dbCopy
}
