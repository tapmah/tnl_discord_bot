package service

import (
	structs "github.com/tapmahtec/TNL_bot"
	"github.com/tapmahtec/TNL_bot/repository"
)

type Activity interface {
	CreateActivity(activ structs.Activities) (int, error)
	GetActivities() ([]structs.Activities, error)
	GetActivityBySid(sid string) (structs.Activities, error)
	DeleteActivityBySid(sid string) error
	AddPlayerActivity(player structs.Players, activity structs.Activities) (int, error)
}

type Players interface {
	AddPlayer(player structs.Players) (int, error)
	GetTopPlayers(limit int) ([]structs.Players, error)
	GetPlayerByName(name string) (structs.Players, error)
	UpdatePlayerScore(id int, score int) error
}

type Service struct {
	Activity
	Players
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Activity: NewActivityService(repo.Activity),
		Players:  NewPlayersService(repo.Players),
	}
}
