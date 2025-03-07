package database

import "fmt"

func (db *appdbimpl) SetMyUserName(username string, cs int) error {
	var righe int64

	err := db.c.QueryRow("SELECT count(username) as righe FROM utenti WHERE id=$1", cs).Scan(&righe)
	if err != nil {
		return fmt.Errorf("user: error find cs in database: %w", err)
	}
	if righe == 0 {
		return fmt.Errorf("non e' stato trovato alcun utente con lo stesso id specificato")
	}

	err = db.c.QueryRow("SELECT count(username) as righe FROM utenti WHERE username= $1", username).Scan(&righe)

	if err != nil {
		return fmt.Errorf("user: error find username in database: %w", err)
	}
	if righe > 0 {
		return fmt.Errorf("l'username scelto e' gia stato utilizzato, sceglierne un altro")
	}

	_, err = db.c.Exec("UPDATE utenti SET username=$1 WHERE id=$2", username, cs)
	if err != nil {
		return fmt.Errorf("error: database UPDATE not successful")
	}
	return nil
}
