package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	db *sql.DB
}

func NewClient(db_path string) (Client, error) {
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return Client{}, err
	}
	c := Client{db}
	err = c.autoMigrate()
	if err != nil {
		return Client{}, err
	}
	return c, nil
}

func (c *Client) autoMigrate() error {

	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);
	`
	_, err := c.db.Exec(usersTable)
	if err != nil {
		return err
	}

	leaguesTable := `
	CREATE TABLE IF NOT EXISTS leagues (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		title TEXT NOT NULL,
		description TEXT,
		database_name TEXT NOT NULL
	);
	`
	_, err = c.db.Exec(leaguesTable)
	if err != nil {
		return err
	}

	leagues_usersTable := `
	CREATE TABLE IF NOT EXISTS users_leagues (
		user_id INTEGER NOT NULL,
		league_id INTEGER NOT NULL,
		UNIQUE("user_id","league_id")
		);
	`
	_, err = c.db.Exec(leagues_usersTable)
	if err != nil {
		return err
	}

	return nil
}
