package coldstore

import (
	"fmt"
	"strconv"
	"strings"
	"thermalkv/internal/model"
	"time"
)

func (m *Manager) MoveToCool(key string, item model.Item) error {

	m.ColdStore.RLock()
	defer m.ColdStore.RUnlock()

	iceCube := fmt.Sprintf("%s|%s|%d\n", key, item.Value, item.Expiry)

	offset, err := m.IceTray.StoreCube(iceCube)

	if err != nil {
		return err
	}

	m.ColdIndex.AddColdEntry(key, offset, item.Expiry)

	return nil
}

func (m *Manager) LoadFromCool(key string) (model.Item, bool) {

	m.ColdStore.RLock()
	defer m.ColdStore.RUnlock()

	offset, exists := m.ColdIndex.MeltColdEntry(key)

	if !exists {
		return model.Item{}, false
	}

	cube := m.IceTray.MeltCube(offset)

	cube = strings.TrimSpace(cube)
	parts := strings.Split(cube, "|")

	if len(parts) < 3 {
		return model.Item{}, false
	}

	expiry, _ := strconv.ParseInt(parts[2], 10, 64)

	return model.Item{
		Value:          parts[1],
		Expiry:         expiry,
		LastAccessUnix: time.Now().Unix(),
		Size:           int64(len(parts[1])),
	}, true

}

func (m *Manager) AppendDelete(
	key string,
) error {

	m.ColdStore.RLock()
	defer m.ColdStore.RUnlock()

	meltedcube := fmt.Sprintf("DEL|%s\n", key)

	_, err := m.IceTray.StoreCube(meltedcube)

	return err
}
