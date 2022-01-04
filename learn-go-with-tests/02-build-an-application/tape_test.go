package poker

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestTape_Write(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()

	tape := &tape{file}

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
