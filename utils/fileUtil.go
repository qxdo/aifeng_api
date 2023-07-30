package utils

import (
	"bufio"
	"os"
)

type File struct {
	file    *os.File
	scanner *bufio.Scanner
}

func OpenFile(path string) (*File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	f := &File{
		file:    file,
		scanner: scanner,
	}
	return f, nil
}

func (f *File) Writeln(content string) error {
	_, err := f.file.WriteString(content)
	return err
}
