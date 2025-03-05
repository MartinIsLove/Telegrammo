package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) DoLogin(username string) (int, error) {
	var if_username string
	var if_id int64
	err := db.c.QueryRow("SELECT username,id FROM utenti WHERE username= $1", username).Scan(&if_username, &if_id)
	if errors.Is(err, sql.ErrNoRows) {
		response, err2 := db.c.Exec("INSERT INTO utenti (username) VALUES ($1)", username)

		if err2 != nil {
			return 0, fmt.Errorf("login: error insert user: %w", err)
		}

		intero, errore := response.LastInsertId()

		if errore != nil {
			return 0, fmt.Errorf("login: error casting insert user: %w", err)
		}

		return int(intero), err2
	}
	if err != nil {
		return 0, fmt.Errorf("login: error query user: %w", err)
	}
	return int(if_id), nil
}
