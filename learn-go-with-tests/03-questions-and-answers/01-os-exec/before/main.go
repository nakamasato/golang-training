package main

import (
	"encoding/xml"
	"log"
	"os/exec"
	"strings"
)

type Payload struct {
	Message string `xml:"message"`
}

func GetData() string {
	cmd := exec.Command("cat", "msg.xml")

	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	var payload Payload
	decoder := xml.NewDecoder(out)

	// these 3 can return errors but I'm ignoring for brevity
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = decoder.Decode(&payload)
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	return strings.ToUpper(payload.Message)
}
