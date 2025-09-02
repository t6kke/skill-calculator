package database

import (
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreateUserParams
}

type CreateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c Client) CreateUser(params CreateUserParams) (User, error) {
	query := `
	INSERT INTO users
	(created_at, updated_at, email, password)
	VALUES
	(CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, ?, ?)
	`
	_, err := c.db.Exec(query, params.Email, params.Password)
	if err != nil {
		return User{}, err
	}

	return c.GetUserByEmail(params.Email)
}

func (c Client) GetUserByEmail(email string) (User, error) {
	query := `
	SELECT id, created_at, updated_at, email, password
	FROM users
	WHERE email = ?
	`
	var user User
	err := c.db.QueryRow(query, email).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, nil
		}
		return User{}, err
	}
	return user, nil
}
