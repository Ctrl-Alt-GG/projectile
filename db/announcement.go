package db

import (
	"sync/atomic"
)

var (
	announcementPtr atomic.Pointer[string]
)

func StoreAnnouncement(newAnnouncement string) error {
	announcementPtr.Store(&newAnnouncement)
	return nil
}

func LoadAnnouncement() string {
	data := announcementPtr.Load()

	if data == nil {
		return ""
	}

	return *data
}
