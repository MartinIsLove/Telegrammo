package database

import "fmt"

func (db *appdbimpl) SetMyUserName(username string, cs int) error {
	var righe int64
	var aut int

	aut, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("user: error in autentication")
	}
	if aut == -1 {
		return err
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

// func (db *appdbimpl) SetMyPhoto(username string, cs int) error {

// }
