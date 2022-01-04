package main

import (
	"fmt"
	"io"
)

type tape struct {
	file io.ReadWriteSeeker
}

func (t *tape) Write(p []byte) (n int, err error) {
	_, err = t.file.Seek(0, 0)
	if err != nil {
		fmt.Println("seek failed")
	}
	return t.file.Write(p)
}
