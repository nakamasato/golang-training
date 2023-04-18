package main

import (
	"fmt"
	"log"

	"ariga.io/atlas/schemahcl"
)

type (
	Family struct {
		Name string `spec:"name,name"`
	}
)

var test struct {
	Families []*Family `spec:"family"`
}

const hcl_str = `
family "default" {
	name = "test"
}
`

func main() {
	fmt.Println("hello")
	// read str into golang struct
	err := schemahcl.New().EvalBytes([]byte(hcl_str), &test, nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(len(test.Families))
	fmt.Println(test.Families[0].Name)
}
