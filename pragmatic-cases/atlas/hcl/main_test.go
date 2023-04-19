package main

import (
	"testing"

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

func Test_EvalBytes(t *testing.T) {
	// read str into golang struct
	err := schemahcl.New().EvalBytes([]byte(hcl_str), &test, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(test.Families) != 1 {
		t.Errorf("Families length want 1, got %d", len(test.Families))
	}
	if want, name := "test", test.Families[0].Name; name != want {
		t.Errorf("Name want %s, got %s", want, name)
	}
}
