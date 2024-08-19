package service

import (
	structs "github.com/tapmahtec/TNL_bot"
	"github.com/tapmahtec/TNL_bot/repository"
)

type PlayersService struct {
	repo repository.Players
}

func NewPlayersService(repo repository.Players) *PlayersService {
	return &PlayersService{repo: repo}
}

func (s *PlayersService) AddPlayer(player structs.Players) (int, error) {
	return s.repo.AddPlayer(player)
}

func (s *PlayersService) GetTopPlayers(limit int) ([]structs.Players, error) {
	return s.repo.GetTopPlayers(limit)
}

func (s *PlayersService) GetPlayerByName(name string) (structs.Players, error) {
	return s.repo.GetPlayerByName(name)
}

func (s *PlayersService) UpdatePlayerScore(id int, score int) error {
	return s.repo.UpdatePlayerScore(id, score)
}
