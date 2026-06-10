package coldstore

import (
	"sync"
	"thermalkv/internal/coldstore/icetray"
	"thermalkv/internal/coldstore/index"
)

type Manager struct {
	ColdIndex *index.ColdIndex
	IceTray   *icetray.IceTray
	ColdStore sync.RWMutex
}

func NewManager() *Manager {
	newtray := icetray.NewIceTray("data/cold.dat")
	newcoldindex := index.NewColdIndex()

	return &Manager{
		ColdIndex: newcoldindex,
		IceTray:   newtray,
	}
}
