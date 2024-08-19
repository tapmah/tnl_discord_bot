package repository

import (
	"database/sql"
	"time"

	structs "github.com/tapmahtec/TNL_bot"
)

type activityMysql struct {
	db *sql.DB
}

func NewActivityMysql(db *sql.DB) *activityMysql {
	return &activityMysql{db: db}
}

func (a *activityMysql) CreateActivity(activ structs.Activities) (int, error) {
	query := `INSERT INTO activity (id, name, sid, score) VALUES (?,?,?,?)`
	result, err := a.db.Exec(query, activ.Id, activ.Name, activ.Sid, activ.Score)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (a *activityMysql) GetActivities() ([]structs.Activities, error) {
	query := `SELECT id, name, sid, score FROM activity ORDER BY score ASC`
	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var activities []structs.Activities
	for rows.Next() {
		var activ structs.Activities
		err := rows.Scan(
			&activ.Id,
			&activ.Name,
			&activ.Sid,
			&activ.Score)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activ)
	}
	return activities, nil
}

func (a *activityMysql) GetActivityBySid(sid string) (structs.Activities, error) {
	query := `SELECT id, name, sid, score FROM activity WHERE sid = ?`
	row := a.db.QueryRow(query, sid)
	var activ structs.Activities
	err := row.Scan(
		&activ.Id,
		&activ.Name,
		&activ.Sid,
		&activ.Score)
	if err == sql.ErrNoRows {
		return structs.Activities{}, nil
	} else if err != nil {
		return structs.Activities{}, err
	}
	return activ, nil
}

func (a *activityMysql) DeleteActivityBySid(sid string) error {
	query := `DELETE FROM activity WHERE sid = ?`
	_, err := a.db.Exec(query, sid)
	return err
}

func (a *activityMysql) AddPlayerActivity(player structs.Players, activity structs.Activities) (int, error) {
	query := `INSERT INTO rel_activity_players (activity_id, player_id, score, times) VALUES (?,?,?,?)`
	result, err := a.db.Exec(query, activity.Id, player.Id, activity.Score, int(time.Now().Unix()))
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}
