package icetray

import "bufio"

func (i *IceTray) MeltCube(offset int64) string {

	i.Mutex.Lock()
	defer i.Mutex.Unlock()

	_, err := i.File.Seek(offset, 0)
	if err != nil {
		return ""
	}

	reader := bufio.NewReader(i.File)
	cube, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}

	return cube
}
