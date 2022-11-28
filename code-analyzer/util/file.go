package util

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

type File struct {
	name   string
	buffer *bytes.Buffer
}

func NewFile(name string) *File {
	return &File{
		name:   name,
		buffer: &bytes.Buffer{},
	}
}

func (f *File) GetName() string {
	return f.name
}

func (f *File) Write(chunk []byte) error {
	_, err := f.buffer.Write(chunk)

	return err
}

func (f *File) ScanAndTmpSave() error {
	path := "./source/" + f.name
	fil, err := os.Create(path)
	if err != nil {
		return err

	}
	buf := make([]byte, 1024)
	for i := 0; i < f.buffer.Len(); i++ {
		// read a chunk
		n, err := f.buffer.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}
		if _, err := fil.Write(buf[:n]); err != nil {
			return err
		}
	}

	fil.Sync()
	defer fil.Close()

	return nil
}

func (f *File) Scan() error {
	scanner := bufio.NewScanner(f.buffer)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
