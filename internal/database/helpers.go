package database

import (
	"database/sql"
	"errors"
	"math/rand"
	"time"
)

func createUnuiqueDBName(c *Client) (string, error) {
	const str_length_length = 16
	const number_of_candidates = 20

	list_of_candidates := make([]string, number_of_candidates)

	for i := range number_of_candidates {
		new_string := generateString(str_length_length)
		list_of_candidates[i] = new_string
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

func generateString(str_length_length int) string {
	rand.Seed(time.Now().UnixNano())
	var character_set = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	result := make([]byte, str_length_length)
	for i := range str_length_length {
		result[i] = character_set[rand.Intn(len(character_set))]
	}
	return string(result)
}
