package db

import (
	"sync"
	"time"

	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
)

var ( // this is a lame, locking array
	gameServers sync.Map
)

func StoreUpdate(id string, update model.GameServer) {
	update.LastUpdate = time.Now()
	gameServers.Store(id, update)
}

func GetServer(id string) (model.GameServer, bool) {
	val, ok := gameServers.Load(id)
	if !ok {
		return model.GameServer{}, false
	}
	srv, ok := val.(model.GameServer)
	if !ok {
		return model.GameServer{}, false
	}

	return srv.Copy(), true
}

func DeleteServer(id string) int {
	_, loaded := gameServers.LoadAndDelete(id)
	if loaded {
		return 1
	}
	return 0
}

// CleanupJob deletes keys that were updated more than limit duration.
// NOTE that this function is not doing its thing atomically, so
// it is possible that when a entry is just updated it's still deleted
func CleanupJob(limit time.Duration) int {
	now := time.Now()
	cnt := 0
	gameServers.Range(func(key, value any) bool {
		v, ok := value.(model.GameServer)
		if !ok {
			cnt++
			gameServers.Delete(key)
			return true
		}

		if now.Sub(v.LastUpdate) > limit {
			cnt++
			gameServers.Delete(key) // apparently you are safe to do this
		}

		return true
	})

	return cnt
}

// GetList makes a deep-ish copy, so later updates in the db does not modify the values returned.
// The keys are not included in GetList
func GetList() []model.GameServer {
	dbCopy := make([]model.GameServer, 0)

	gameServers.Range(func(_, value any) bool {
		v, ok := value.(model.GameServer)
		if ok {
			dbCopy = append(dbCopy, v.Copy())
		}

		return true // continue
	})

	return dbCopy
}

// GetMap is similar to GetList, but it returns a map instead, so keys are "visible".
func GetMap() map[string]model.GameServer {
	dbCopy := make(map[string]model.GameServer)

	gameServers.Range(func(key, value any) bool {
		k, ok := key.(string)
		if !ok {
			return true
		}
		v, ok := value.(model.GameServer)
		if !ok {
			return true
		}

		dbCopy[k] = v.Copy()
		return true // continue
	})

	return dbCopy
}
