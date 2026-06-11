package index

import (
	"sync"
)

type ColdEntry struct {
	Offset int64
	Expiry int64
}

type ColdIndex struct {
	entries map[string]ColdEntry
	Mutex   sync.RWMutex
}

func NewColdIndex() *ColdIndex {

	return &ColdIndex{
		entries: make(map[string]ColdEntry),
	}

}

func (i *ColdIndex) AddColdEntry(key string, offset, expiry int64) {
	i.Mutex.Lock()
	defer i.Mutex.Unlock()

	i.entries[key] = ColdEntry{
		Offset: offset,
		Expiry: expiry,
	}
}

func (i *ColdIndex) MeltColdEntryNoLock(key string) (int64, bool) {

	entry, exists := i.entries[key]

	if !exists {
		return 0, false
	}

	return entry.Offset, true

}

func (i *ColdIndex) MeltColdEntry(key string) (int64, bool) {

	i.Mutex.RLock()
	defer i.Mutex.RUnlock()
	return i.MeltColdEntryNoLock(key)

}

func (i *ColdIndex) DestroyColdEntry(key string) {
	i.Mutex.Lock()
	defer i.Mutex.Unlock()

	delete(i.entries, key)
}

func (i *ColdIndex) HaveColdEntry(key string) bool {
	_, exists := i.entries[key]
	return exists
}

func (i *ColdIndex) GetSize() int {
	return len(i.entries)
}

func (i *ColdIndex) Keys() []string {
	i.Mutex.RLock()
	defer i.Mutex.RUnlock()

	keys := make([]string, 0, len(i.entries))

	for k := range i.entries {
		keys = append(keys, k)
	}

	return keys
}

func (i *ColdIndex) ReplaceEntries(newEntries map[string]ColdEntry) {
	i.Mutex.Lock()
	defer i.Mutex.Unlock()

	i.entries = newEntries
}
