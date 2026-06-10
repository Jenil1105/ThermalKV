package icetray

import (
	"fmt"
	"os"
)

func (i *IceTray) ChangeIceTray(path string) {
	i.Mutex.Lock()
	defer i.Mutex.Unlock()

	err := i.File.Close()
	if err != nil {
		return
	}

	original := "data/cold.dat"

	err = os.Rename(path, original)
	if err != nil {
		return
	}

	file, err := os.OpenFile(
		original,
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		0644,
	)

	if err != nil {
		fmt.Println("Ice Tray open error: ", err)
	}

	i.File = file
	fmt.Println("Compaction Successful")
}
