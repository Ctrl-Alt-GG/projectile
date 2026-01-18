package db

import (
	"slices"
	"sync"
	"time"

	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
)

var ( // this is a lame, locking array
	gameServers     []model.GameServer
	gameServersLock sync.RWMutex
)

func StoreUpdate(update model.GameServer) error {
	gameServersLock.Lock()
	defer gameServersLock.Unlock()

	update.LastUpdate = time.Now()

	idx := slices.IndexFunc(gameServers, func(server model.GameServer) bool {
		return server.ID.Equals(update.ID)
	})

	if idx == -1 {
		gameServers = append(gameServers, update)
	} else {
		gameServers[idx] = update
	}

	return nil
}

func DeleteServer(id model.Identifier) (int, error) {
	gameServersLock.Lock()
	defer gameServersLock.Unlock()

	lenBefore := len(gameServers)

	gameServers = slices.DeleteFunc(gameServers, func(server model.GameServer) bool {
		return server.ID.Equals(id)
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

// GetAll makes a deep-ish copy, so later updates in the db does not modify the values returned.
func GetAll() []model.GameServer {
	gameServersLock.RLock()
	defer gameServersLock.RUnlock()

	dbCopy := make([]model.GameServer, len(gameServers))

	for i, val := range gameServers {
		dbCopy[i] = val.Copy()
	}

	return dbCopy
}
