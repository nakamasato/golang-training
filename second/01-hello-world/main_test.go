package main

import (
	"testing"
)

func TestHello(t *testing.T)  {
	assertCorrectMessage := func(t *testing.T, got, want string) {
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	}
	t.Run("Hello with naka", func(t *testing.T) {
		got := Hello("Naka", "")
		want := "Hello, Naka!"
		assertCorrectMessage(t, got, want)
	})
	t.Run("Hello with empty string", func(t *testing.T) {
		got := Hello("", "")
		want := "Hello, World!"
		assertCorrectMessage(t, got, want)
	})

	t.Run("In Spanish", func(t *testing.T) {
		got := Hello("Elodie", "Spanish")
		want := "Ola, Elodie!"
		assertCorrectMessage(t, got, want)
	})
}
