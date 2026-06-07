package recover

import (
	"thermalkv/internal/persistence/snapshot"
	"thermalkv/internal/store"
)

func RecoverSnapshot(s *store.Store, path string) {
	snaps := snapshot.LoadSnapshot(path)
	s.ImportData(snaps)
}
