package repository

import (
	"database/sql"

	structs "github.com/tapmahtec/TNL_bot"
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

type Repository struct {
	Activity
	Players
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Activity: NewActivityMysql(db),
		Players:  NewPlayersMysql(db),
	}
}
