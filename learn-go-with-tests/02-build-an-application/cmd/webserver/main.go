package main

import (
	"log"
	"net/http"
	"os"
	"tmp/learn-go-with-tests/02-build-an-application"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := poker.NewFileSystemPlayerStore(db)
	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}
	server := poker.NewPlayerServer(store)
	// server := NewPlayerServer(NewInMemoryPlayerStore())
	log.Fatal(http.ListenAndServe(":5000", server))
}
