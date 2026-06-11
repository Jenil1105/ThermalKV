package store

import (
	"fmt"
	"strconv"
	"time"
)

// Set TTL
func (s *Store) SetTTL(key string, seconds int) {

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	expiry := time.Now().Unix() + int64(seconds)
	isExpirySet := s.setItemExpiry(key, expiry)
	if !isExpirySet {
		return
	}

	err := s.WAL.Write("EXPIRE", key, strconv.FormatInt(expiry, 10))
	if err != nil {
		fmt.Println("WAL write failed: ", err)
		return
	}
}
