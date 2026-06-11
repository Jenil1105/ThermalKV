package coldstore

func (m *Manager) AddIndex(key string, offset, expiry int64) {
	m.ColdIndex.AddColdEntry(key, offset, expiry)
}

func (m *Manager) DeleteIndex(key string) {
	m.ColdIndex.DestroyColdEntry(key)
}

func (m *Manager) HaveIndex(key string) bool {
	return m.ColdIndex.HaveColdEntry(key)
}

func (m *Manager) GetIndexSize() int {
	return m.ColdIndex.GetSize()
}
