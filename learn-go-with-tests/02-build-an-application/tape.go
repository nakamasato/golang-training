package poker

import (
	"fmt"
	"os"
)

type Tape struct {
	File *os.File
}

func (t *Tape) Write(p []byte) (n int, err error) {
	err = t.File.Truncate(0)
	if err != nil {
		fmt.Println("truncate failed")
	}
	_, err = t.File.Seek(0, 0)
	if err != nil {
		fmt.Println("seek failed")
	}
	return t.File.Write(p)
}
