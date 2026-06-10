package icetray

import "io"

func (i *IceTray) StoreCube(iceCube string) (int64, error) {

	i.Mutex.Lock()
	defer i.Mutex.Unlock()

	offset, err := i.File.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	_, err = i.File.WriteString(iceCube)

	if err != nil {
		return 0, err
	}

	return offset, nil
}
