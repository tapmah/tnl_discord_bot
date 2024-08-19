package service

import (
	"errors"

	structs "github.com/tapmahtec/TNL_bot"
	"github.com/tapmahtec/TNL_bot/repository"
)

type ActivityService struct {
	repo repository.Activity
}

func NewActivityService(repo repository.Activity) *ActivityService {
	return &ActivityService{repo: repo}
}

func (a *ActivityService) CreateActivity(activ structs.Activities) (int, error) {
	if len(activ.Name) == 0 {
		return 0, errors.New("имя не может быть пустым")
	}
	if len(activ.Sid) == 0 {
		return 0, errors.New("системное имя не может быть пустым")
	}
	if activ.Score == 0 {
		return 0, errors.New("счет не может быть пустым")
	}

	return a.repo.CreateActivity(activ)
}

func (a *ActivityService) GetActivities() ([]structs.Activities, error) {
	return a.repo.GetActivities()
}

func (a *ActivityService) GetActivityBySid(sid string) (structs.Activities, error) {
	if len(sid) == 0 {
		return structs.Activities{}, errors.New("системное имя не может быть пустым")
	}
	return a.repo.GetActivityBySid(sid)
}

func (a *ActivityService) DeleteActivityBySid(sid string) error {
	if len(sid) == 0 {
		return errors.New("системное имя не может быть пустым")
	}
	act, err := a.repo.GetActivityBySid(sid)
	if err != nil {
		return err
	}
	if act.Id == 0 {
		return errors.New("не найдена указанная активность")
	}
	return a.repo.DeleteActivityBySid(sid)
}

func (a *ActivityService) AddPlayerActivity(player structs.Players, activity structs.Activities) (int, error) {
	return a.repo.AddPlayerActivity(player, activity)
}
