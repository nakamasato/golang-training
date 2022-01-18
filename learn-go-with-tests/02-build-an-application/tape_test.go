package poker_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"tmp/learn-go-with-tests/02-build-an-application"
)

func TestTape_Write(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()

	tape := &poker.Tape{file}

	_, err := tape.Write([]byte("abc"))
	if err != nil {
		fmt.Println("failed to Write")
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println("failed to Seek")
	}
	newFileContents, _ := ioutil.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
