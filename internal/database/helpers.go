package database

import (
	"crypto/rand"
	"database/sql"
	"errors"
)

func createUnuiqueDBName(c *Client) (string, error) {
	const number_of_candidates = 20

	list_of_candidates := make([]string, number_of_candidates)

	for i := range number_of_candidates {
		list_of_candidates[i] = rand.Text()
	}

	query := `
	SELECT id
	FROM leagues
	WHERE database_name = ?
	`
	var id int
	for i := range len(list_of_candidates) {
		err := c.db.QueryRow(query, list_of_candidates[i]).Scan(&id)
		if errors.Is(err, sql.ErrNoRows) {
			return list_of_candidates[i], nil
		}
		if err != nil {
			return "", nil
		}
	}

	return "", errors.New("failed to generate unique database string")
}
