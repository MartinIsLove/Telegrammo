package database

import "fmt"

func (db *appdbimpl) Authentication(cs int) (int, error) {
	var righe int64
	err := db.c.QueryRow("SELECT count(username) as righe FROM utenti WHERE id=$1", cs).Scan(&righe)
	if err != nil {
		return -1, fmt.Errorf("user: error find cs in database: %w", err)
	}
	if righe == 0 {
		return -1, fmt.Errorf("non e' stato trovato alcun utente con lo stesso id specificato")
	}
	return 1, nil
}
