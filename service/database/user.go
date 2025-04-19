package database

import "fmt"

// imposta lo username dell'utente loggato
func (db *appdbimpl) SetMyUserName(username string, cs int) error {
	var righe int64

	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("SetMyUserName: error in authentication: %w", err)
	}

	// conta quanti username esistono nel database uguali a quello inserito in input nel quale vuole cambiare il proprio username l'utente loggato
	err = db.c.QueryRow("SELECT count(username) as righe FROM utenti WHERE username= $1", username).Scan(&righe)

	if err != nil {
		return fmt.Errorf("SetMyUserName: error find username in database: %w", err)
	}

	// se sono piu' di 0 gli username uguali a quello scelto
	if righe > 0 {
		return fmt.Errorf("username already used, choose another one: %w", err)
	}

	// cambia l'username dell'utente loggato
	_, err = db.c.Exec("UPDATE utenti SET username=$1 WHERE id=$2", username, cs)
	if err != nil {
		return fmt.Errorf("SetMyUserName: database UPDATE not successful: %w", err)
	}
	return nil
}

// imposta la foto dell'utente loggato
func (db *appdbimpl) SetMyPhoto(photo []byte, cs int) error {

	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("SetMyPhoto: error in authentication: %w", err)
	}

	// cambia la foto dell'utente loggato
	_, err = db.c.Exec("UPDATE utenti SET propic=$1 WHERE id=$2", photo, cs)

	if err != nil {
		return fmt.Errorf("SetMyPhoto: error in UPDATE database")
	}
	return nil

}

// restituisce i parametri dell'utente loggato
func (db *appdbimpl) GetMyUser(cs int) (string, []byte, int, error) {

	var username string
	var propic []byte
	var id int

	_, err := db.Authentication(cs)
	if err != nil {
		var tmp []byte
		return "", tmp, 0, fmt.Errorf("GetMyUser: error in authentication: %w", err)
	}

	// seleziono i parametri dell'utente loggato
	err = db.c.QueryRow("SELECT id, username, propic FROM utenti WHERE id=$1", cs).Scan(&id, &username, &propic)
	if err != nil {
		var tmp []byte
		return "", tmp, 0, fmt.Errorf("GetMyUser: error find username in database: %w", err)
	}

	return username, propic, id, nil
}

// ritorna true se l'utente esiste false se non esiste
func (db *appdbimpl) UserExist(id int) bool {
	var tmp int

	// conta il numero di occorrenze di id dove l'id e' uguale all'id cercato
	err := db.c.QueryRow("SELECT count(id) FROM utenti WHERE id=$1;", id).Scan(&tmp)

	if err != nil {
		return false
	}
	if tmp == 0 {
		return false
	}

	return true
}
