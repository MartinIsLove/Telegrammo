package database

import "fmt"

// funzione di autenticazione dell'utente loggato restituisce 1 se l'utente esiste nel database 0 altrimenti
func (db *appdbimpl) Authentication(cs int) (int, error) {
	var righe int64

	// conta quanti utenti con l'id in input sono presenti nel database
	err := db.c.QueryRow("SELECT count(username) as righe FROM utenti WHERE id=$1", cs).Scan(&righe)
	if err != nil {
		return -1, fmt.Errorf("authentication: query error find cs in database: %w", err)
	}
	if righe == 0 {
		return -1, fmt.Errorf("don't find any user with this id")
	}
	return 1, nil
}
