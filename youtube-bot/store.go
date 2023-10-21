package main

import "sync"

const storeSize = 128 // How many user data can store at once

type VideoDataStore struct {
	data map[string][]VideoData
	mu   sync.Mutex
}

var videoStore = &VideoDataStore{
	data: make(map[string][]VideoData, storeSize),
}

func checkStoreSize(currentUserId string) {
	if len(videoStore.data) < storeSize {
		return
	}

	for key := range videoStore.data {
		if key != currentUserId {
			delete(videoStore.data, key)
		}
	}
}
