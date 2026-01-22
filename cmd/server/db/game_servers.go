package db

import (
	"sync"
	"time"

	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
)

var ( // this is a lame, locking array
	gameServers sync.Map
)

func StoreGameServerData(id string, update model.GameServerData) {
	sgsd := model.StoredGameServerData{
		GameServerData: update,
		LastUpdate:     time.Now(),
	}
	gameServers.Store(id, sgsd)
}

func GetServer(id string) (model.StoredGameServerData, bool) {
	val, ok := gameServers.Load(id)
	if !ok {
		return model.StoredGameServerData{}, false
	}
	srv, ok := val.(model.StoredGameServerData)
	if !ok {
		return model.StoredGameServerData{}, false
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
// it is possible that when an entry is just updated it's still deleted
func CleanupJob(limit time.Duration) int {
	now := time.Now()
	cnt := 0
	gameServers.Range(func(key, value any) bool {
		v, ok := value.(model.StoredGameServerData)
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
// The keys and update timestamps are not included in GetList
func GetList() []model.GameServerData {
	dbCopy := make([]model.GameServerData, 0)

	gameServers.Range(func(_, value any) bool {
		v, ok := value.(model.StoredGameServerData)
		if ok {
			dbCopy = append(dbCopy, v.GameServerData.Copy())
		}

		return true // continue
	})

	return dbCopy
}

// GetMap is similar to GetList, but it returns a map instead, so keys are "visible".
func GetMap() map[string]model.StoredGameServerData {
	dbCopy := make(map[string]model.StoredGameServerData)

	gameServers.Range(func(key, value any) bool {
		k, ok := key.(string)
		if !ok {
			return true
		}
		v, ok := value.(model.StoredGameServerData)
		if !ok {
			return true
		}

		dbCopy[k] = v.Copy()
		return true // continue
	})

	return dbCopy
}
