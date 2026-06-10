package store

import (
	"fmt"
	"os"
)

func (s *Store) GetInfo() []string {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	info := []string{
		"===== ThermalKV Info =====",
		fmt.Sprintf(
			"HOT Keys           : %d",
			len(s.Data),
		),
		fmt.Sprintf(
			"COOL Keys          : %d",
			s.Thermal.GetIndexSize(),
		),
		fmt.Sprintf(
			"HOT Memory Usage   : %d bytes",
			s.CurrentMemoryUsage,
		),
		fmt.Sprintf(
			"Max HOT Memory     : %d bytes",
			s.MaxHotMemory,
		),
		fmt.Sprintf(
			"Cooling Threshold  : %d",
			s.CoolingThreshold,
		),
		fmt.Sprintf(
			"Cold File Size     : %d",
			GetColdFileSize(),
		),
	}

	return info

}

func GetColdFileSize() int64 {
	fileInfo, err := os.Stat("data/cold.dat")

	if err != nil {
		return 0
	}
	return fileInfo.Size()
}
