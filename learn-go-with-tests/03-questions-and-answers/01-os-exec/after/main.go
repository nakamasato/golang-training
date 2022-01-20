package main

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

type Payload struct {
	Message string `xml:"message"`
}

func GetData(data io.Reader) string {
	var payload Payload
	xml.NewDecoder(data).Decode(&payload)
	return strings.ToUpper(payload.Message)
}

func getXMLFromCommand() io.Reader {
	cmd := exec.Command("cat", "msg.xml")
	out, _ := cmd.StdoutPipe()

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(out)
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	return bytes.NewReader(data)
}
