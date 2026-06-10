package icetray

import (
	"fmt"
	"os"
	"sync"
)

type IceTray struct {
	File  *os.File
	Mutex sync.RWMutex
}

func NewIceTray(path string) *IceTray {

	file, err := os.OpenFile(
		path,
		os.O_RDWR|os.O_APPEND|os.O_CREATE,
		0644,
	)

	if err != nil {
		fmt.Println("Ice Tray open error: ", err)
		return nil
	}

	return &IceTray{
		File: file,
	}
}
