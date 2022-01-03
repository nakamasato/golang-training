package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

type Message struct {
	Name, Text string
}

func main() {
	const jsonStream = `
	[
		{"Name": "Ed", "Text": "Knock knock."},
		{"Name": "Sam", "Text": "Who's there?"},
		{"Name": "Ed", "Text": "Go fmt."},
		{"Name": "Sam", "Text": "Go fmt who?"},
		{"Name": "Ed", "Text": "Go fmt yourself!"}]`

	messageWithStringsReader := decodeWithStringsReader(jsonStream)

	messageWithOsFile := decodeWithOsFile(jsonStream)

	if !reflect.DeepEqual(messageWithStringsReader, messageWithStringsReader) {
		fmt.Printf("Results are different strings.NewReader: %v, io.File: %v\n", messageWithStringsReader, messageWithOsFile)
	}
	fmt.Println("Results are same!")
}

func decodeWithStringsReader(jsonStream string) []Message {
	var messages []Message
	err := json.NewDecoder(strings.NewReader(jsonStream)).Decode(&messages)
	if err != nil {
		fmt.Printf("Failed to decode with strings.Reader %v", err)
		return nil
	}
	return messages
}

func decodeWithOsFile(jsonStream string) []Message {
	var messages []Message
	tmpfile, err := ioutil.TempFile(".", "db")
	if err != nil {
		fmt.Println("Failed to create tmpfile")
		return nil
	}
	tmpfile.Write([]byte(jsonStream))
	fmt.Println(tmpfile.Name())
	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}
	defer removeFile()

	tmpfile.Seek(0, 0) // read from the beginning
	err = json.NewDecoder(tmpfile).Decode(&messages)
	if err != nil {
		fmt.Printf("Failed to decode with os.File %v", err)
		return nil
	}
	return messages
}
