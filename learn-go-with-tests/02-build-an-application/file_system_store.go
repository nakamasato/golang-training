package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type FileSystemPlayerStore struct {
	database *json.Encoder
	league   League
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	_, err := file.Seek(0, 0)
	if err != nil {
		fmt.Println("seek failed")
	}

	// To support empty file with content ''
	err = initialisePlayerDBFile(file)
	if err != nil {
		fmt.Println("failed to initialize player db")
	}

	league, err := NewLeague(file)
	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}
	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil
}

func initialisePlayerDBFile(file *os.File) error {
	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		_, err = file.Write([]byte("[]"))
		if err != nil {
			fmt.Println("Failed to write")
		}
		_, err = file.Seek(0, 0)
		if err != nil {
			fmt.Printf("Failed to seek")
		}
	}
	return nil
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
	err := f.database.Encode(f.league)
	if err != nil {
		fmt.Println("Encode failed")
	}
}
