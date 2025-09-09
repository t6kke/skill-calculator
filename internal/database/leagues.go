package database

import (
	"database/sql"
	"errors"
	"time"
)

type League struct {
	ID            int       `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreateLeagueParams
	DatabaseName string     `json:"database_name"`
}

type CreateLeagueParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      int    `json:"user_id"`
}

func (c Client) CreateLeageWithUserRelation(params CreateLeagueParams) (League, error) {
	db_name, err := createUnuiqueDBName(&c)
	if err != nil {
		return League{}, err
	}
	db_name = db_name + ".db"

	create_league_query := `
	INSERT INTO leagues
	(created_at, updated_at, title, description, database_name)
	VALUES
	(CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, ?, ?, ?)
	`
	create_user_leage_relation_query := `
	INSERT INTO users_leagues
	(user_id, league_id)
	VALUES
	(?, ?)
	`

	result, err := c.db.Exec(create_league_query, params.Title, params.Description, db_name)
	if err != nil {
		return League{}, err
	}
	league_id, _ := result.LastInsertId() //TODO do error hanlding?

	_, err = c.db.Exec(create_user_leage_relation_query, params.UserID, league_id)
	if err != nil {
		//TODO the entry for the league creation should be deleted
		return League{}, err
	}

	league, err := c.GetLeague(int(league_id))
	if err != nil {
		return League{}, err
	}

	return league, nil
}

func (c Client) GetLeague(league_id int) (League, error) {
	query := `
	SELECT l.id, l.created_at, l.updated_at, l.title, l.description, l.database_name, ul.user_id
	FROM leagues l
	JOIN users_leagues ul on ul.league_id = l.id
	WHERE l.id = ?
	`

	var league League
	err := c.db.QueryRow(query, league_id).Scan(
		&league.ID,
		&league.CreatedAt,
		&league.UpdatedAt,
		&league.Title,
		&league.Description,
		&league.DatabaseName,
		&league.UserID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return League{}, nil
		}
		return League{}, err
	}

	return league, err
}

func (c Client) GetLeagues(user_id int) ([]League, error) {
	query := `
	SELECT l.id, l.created_at, l.updated_at, l.title, l.description, l.database_name, ul.user_id
	FROM leagues l
	JOIN users_leagues ul on ul.league_id = l.id
	WHERE ul.user_id = ?
	`
	rows, err := c.db.Query(query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	leagues := []League{}

	for rows.Next() {
		var league League
		if err := rows.Scan(
			&league.ID,
			&league.CreatedAt,
		        &league.UpdatedAt,
		        &league.Title,
		        &league.Description,
		        &league.DatabaseName,
		        &league.UserID,
		); err != nil {
			return nil, err
		}
		leagues = append(leagues, league)
	}

	return leagues, nil
}
