package icetray

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type IceTray struct {
	File   *os.File
	Writer *bufio.Writer
}

func NewIceTray(path string) *IceTray {

	file, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		fmt.Println("Ice Tray open error: ", err)
		return nil
	}

	writer := bufio.NewWriterSize(file, 64*1024)

	return &IceTray{
		File:   file,
		Writer: writer,
	}
}

func (i *IceTray) StoreCube(iceCube string) (int64, error) {

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
