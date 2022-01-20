package main

import "testing"

func TestGetData(t *testing.T) {
	got := GetData()
	want := "HAPPY NEW YEAR!"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
