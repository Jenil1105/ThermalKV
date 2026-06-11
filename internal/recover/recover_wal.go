package recover

import (
	"strconv"
	"strings"
	"thermalkv/internal/persistence/walpkg"
	"thermalkv/internal/store"
)

// Recover replays the given logs to restore the store's state after a restart.
func StoreLogs(s *store.Store, logs []string) {
	for _, log := range logs {
		parts := strings.Fields(log)

		if len(parts) < 2 {
			continue
		}

		operation := parts[0]
		key := parts[1]

		switch operation {

		case "SET":
			if len(parts) < 3 {
				continue
			}

			value := strings.Join(parts[2:], " ")
			s.RecoverSet(key, value)

		case "DEL":

			s.RecoverDelete(key)

		case "EXPIRE":

			if len(parts) < 3 {
				continue
			}

			expiry, err := strconv.ParseInt(parts[2], 10, 64)
			if err != nil {
				continue
			}
			s.RecoverExpire(key, expiry)
		}
	}
}

func RecoverWAL(s *store.Store, dir string) {

	logs := walpkg.LoadLogs(dir)
	StoreLogs(s, logs)

}
