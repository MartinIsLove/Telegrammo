package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// effettua il login dell'utente
func (db *appdbimpl) DoLogin(username string) (int, error) {
	var if_username string
	var if_id int64

	// seleziona l'utente con lo username inserito dalla tabella utenti e se esiste ritorna il suo id
	err := db.c.QueryRow("SELECT username,id FROM utenti WHERE username= $1", username).Scan(&if_username, &if_id)

	//se non esiste lo crea
	if errors.Is(err, sql.ErrNoRows) {
		noPhotoPath := filepath.Join("/workspace/webui/src/assets/", "NoPhoto.png")
		noPhotoBytes, err := os.ReadFile(noPhotoPath)
		if err != nil {
			return 0, fmt.Errorf("login: error reading noPhoto.svg: %w", err)
		}

		// inserisce l'utente nel database
		response, err2 := db.c.Exec("INSERT INTO utenti (username, propic) VALUES ($1, $2)", username, noPhotoBytes)

		if err2 != nil {
			return 0, fmt.Errorf("login: error insert user: %w", err)
		}

		// recupero l'id dell'utente appena inserito
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
