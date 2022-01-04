package main

import (
	"fmt"
	"os"
)

type tape struct {
	file *os.File
}

func (t *tape) Write(p []byte) (n int, err error) {
	err = t.file.Truncate(0)
	if err != nil {
		fmt.Println("truncate failed")
	}
	_, err = t.file.Seek(0, 0)
	if err != nil {
		fmt.Println("seek failed")
	}
	return t.file.Write(p)
}
