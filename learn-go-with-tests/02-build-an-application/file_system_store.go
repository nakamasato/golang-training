package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type FileSystemPlayerStore struct {
	database io.Writer
	league   League
}

func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
	_, err := database.Seek(0, 0)
	if err != nil {
		fmt.Println("seek failed")
	}
	league, _ := NewLeague(database)
	return &FileSystemPlayerStore{
		database: &tape{database},
		league:   league,
	}
}

func (f *FileSystemPlayerStore) GetLeague() League {
	return f.league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	player := f.GetLeague().Find(name)
	if player != nil {
		return player.Wins
	}
	return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{Name: name, Wins: 1})
	}
<<<<<<< HEAD

=======
>>>>>>> 4e6b25e ([learn go with tests][02-build-an-application] IO and sorting Step 10: More refactoring and performance concerns)
	err := json.NewEncoder(f.database).Encode(f.league)
	if err != nil {
		fmt.Println("Encode failed")
	}
}
