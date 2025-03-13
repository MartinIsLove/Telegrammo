package database

import "fmt"

func (db *appdbimpl) SetMyUserName(username string, cs int) error {
	var righe int64

	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("user: error in authentication: %w", err)
	}

	err = db.c.QueryRow("SELECT count(username) as righe FROM utenti WHERE username= $1", username).Scan(&righe)

	if err != nil {
		return fmt.Errorf("user: error find username in database: %w", err)
	}
	if righe > 0 {
		return fmt.Errorf("l'username scelto e' gia stato utilizzato, sceglierne un altro: %w", err)
	}

	_, err = db.c.Exec("UPDATE utenti SET username=$1 WHERE id=$2", username, cs)
	if err != nil {
		return fmt.Errorf("error: database UPDATE not successful: %w", err)
	}
	return nil
}

func (db *appdbimpl) SetMyPhoto(photo []byte, cs int) error {

	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("error in authentication: %w", err)
	}
	_, err = db.c.Exec("UPDATE utenti SET propic=$1 WHERE id=$2", photo, cs)

	if err != nil {
		return fmt.Errorf("error in UPDATE db")
	}
	return nil

}

func (db *appdbimpl) GetMyUser(cs int) (string, []byte, int, error) {

	var username string
	var propic []byte
	var id int

	_, err := db.Authentication(cs)
	if err != nil {
		var tmp []byte
		return "", tmp, 0, fmt.Errorf("error in authentication: %w", err)
	}

	err = db.c.QueryRow("SELECT id, username, propic FROM utenti WHERE id=$1", cs).Scan(&id, &username, &propic)
	if err != nil {
		var tmp []byte
		return "", tmp, 0, fmt.Errorf("user: error find username in database: %w", err)
	}

	return username, propic, id, nil
}
