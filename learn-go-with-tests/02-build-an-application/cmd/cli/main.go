package main

import (
	"fmt"
	"log"
	"os"
	"tmp/learn-go-with-tests/02-build-an-application"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}
	poker.NewCLI(store, os.Stdin, poker.BlindAlerterFunc(poker.StdOutAlerter)).PlayPoker()
}
