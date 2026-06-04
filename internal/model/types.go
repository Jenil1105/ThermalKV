package model

type Item struct {
	Value          string
	Expiry         int64
	LastAccessUnix int64
}

type SnapshotItem struct {
	Value          string
	Expiry         int64
	LastAccessUnix int64
}
