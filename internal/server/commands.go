package server

import (
	"fmt"
	"strconv"
	"strings"
)

func ExecuteCommand(db KVService, input string) ([]string, bool) {
	input = strings.TrimSpace(input)
	if input == "" {
		return []string{"Empty command"}, false
	}

	parts := strings.Split(input, " ")
	command := strings.ToUpper(parts[0])

	switch command {
	case "SET":
		if len(parts) < 3 {
			return []string{"Usage: SET key value"}, false
		}

		key := parts[1]
		value := strings.Join(parts[2:], " ")
		db.Set(key, value)
		return []string{"OK :)"}, false

	case "GET":
		if len(parts) < 2 {
			return []string{"Usage: GET key"}, false
		}

		key := parts[1]
		value, exists := db.Get(key)
		if exists {
			return []string{value}, false
		}
		return []string{"Key not found... :("}, false

	case "DEL":
		if len(parts) < 2 {
			return []string{"Usage: DEL key"}, false
		}

		key := parts[1]
		db.Delete(key)
		return []string{"OK :)"}, false

	case "TTL":
		if len(parts) < 3 {
			return []string{"Usage: TTL key seconds"}, false
		}

		key := parts[1]
		seconds, err := strconv.Atoi(parts[2])
		if err != nil {
			return []string{"Invalid seconds :/"}, false
		}

		_, exists := db.Get(key)
		if exists {
			db.SetTTL(key, seconds)
			return []string{"OK :)"}, false
		}
		return []string{"Key not found... :("}, false

	case "COOL":
		if len(parts) < 2 {
			return []string{"Usage: COOL key"}, false
		}

		key := parts[1]
		err := db.CoolKey(key)
		if err != nil {
			return []string{err.Error()}, false
		}
		return []string{"OK :)"}, false

	case "COUNT":
		count := db.Count()
		return []string{fmt.Sprintf("%d", count)}, false

	case "EXISTS":
		if len(parts) < 2 {
			return []string{"Usage: EXISTS key"}, false
		}

		key := parts[1]
		if db.Exists(key) {
			return []string{"true"}, false
		}
		return []string{"false"}, false

	case "KEYS":
		keys := db.Keys()
		if len(keys) == 0 {
			return []string{"No keys :|"}, false
		}
		return []string{strings.Join(keys, ", ")}, false

	case "INFO":
		return db.GetInfo(), false

	case "EXIT":
		return []string{"bye... ;) "}, true

	default:
		return []string{"Unknown command :/"}, false
	}
}
