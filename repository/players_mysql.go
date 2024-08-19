package repository

import (
	"database/sql"

	structs "github.com/tapmahtec/TNL_bot"
)

type PlayersMysql struct {
	db *sql.DB
}

func NewPlayersMysql(db *sql.DB) *PlayersMysql {
	return &PlayersMysql{db: db}
}

func (p *PlayersMysql) AddPlayer(player structs.Players) (int, error) {
	query := `INSERT INTO players (name, score) VALUES (?,?)`
	result, err := p.db.Exec(query, player.Name, player.Score)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (p *PlayersMysql) UpdatePlayerScore(id int, score int) error {
	query := `UPDATE players SET score = ? WHERE id = ?`
	_, err := p.db.Exec(query, score, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PlayersMysql) GetTopPlayers(limit int) ([]structs.Players, error) {
	query := `SELECT id, name, score FROM players ORDER BY score DESC LIMIT ?`
	rows, err := p.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var players []structs.Players
	for rows.Next() {
		var player structs.Players
		err := rows.Scan(&player.Id, &player.Name, &player.Score)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, nil
}

func (p *PlayersMysql) GetPlayerByName(name string) (structs.Players, error) {
	query := `SELECT id, name, score FROM players WHERE name =?`
	row := p.db.QueryRow(query, name)
	var player structs.Players
	err := row.Scan(&player.Id, &player.Name, &player.Score)
	if err == sql.ErrNoRows {
		return structs.Players{}, nil
	} else if err != nil {
		return structs.Players{}, err
	}
	return player, nil
}
