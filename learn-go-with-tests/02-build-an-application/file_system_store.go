package main

import (
	"encoding/json"
	"io"
)

type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
}

func (f FileSystemPlayerStore) GetLeague() League {
	f.database.Seek(0, 0)
	league, err := NewLeague(f.database)
	if err != nil {
		return []Player{}
	}
	return league
}

func (f FileSystemPlayerStore) GetPlayerScore(name string) int {
	player := f.GetLeague().Find(name)
	if player != nil {
		return player.Wins
	}
	return 0
}

func (f FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()
	player := league.Find(name)
	if player != nil {
		player.Wins++
	} else {
		league = append(league, Player{Name: name, Wins: 1})
	}
	f.database.Seek(0, 0)
	json.NewEncoder(f.database).Encode(league)
}
