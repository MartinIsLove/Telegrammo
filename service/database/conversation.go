package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) IsChatDuplicated(cs int, id int) (bool, error) {
	var num_righe int64

	err := db.c.QueryRow("SELECT COUNT(chat.id) AS righe FROM chat JOIN membri ON chat.id=membri.id_chat WHERE chat.gruppo=FALSE AND membri.id_utenti IN ($1 , $2) GROUP BY chat.id HAVING COUNT (DISTINCT membri.id_utenti)=2", cs, id).Scan(&num_righe)
	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return true, fmt.Errorf("chat: error checking chat duplicated in database: %w", err)
	}

	if num_righe == 0 {
		return false, nil
	}
	return true, fmt.Errorf("chat: error this chat already exist: %w", err)
}

func (db *appdbimpl) CreateChat(cs int, id int) error {

	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("error in authentication CreateChat: %w", err)
	}

	Isduplicated, err := db.IsChatDuplicated(cs, id)

	if err != nil {
		return err
	}

	if !Isduplicated {
		str_chat, err := db.c.Exec("INSERT INTO chat (gruppo) VALUES (FALSE)")
		if err != nil {
			return fmt.Errorf("chat: error insert chat in table chat: %w", err)
		}
		id_chat, err := str_chat.LastInsertId()
		if err != nil {
			return fmt.Errorf("chat: error catch number of rows from query: %w", err)
		}

		_, err = db.c.Exec("INSERT INTO membri (id_chat, id_utenti) VALUES  ($1, $2), ($1, $3)", id_chat, cs, id)
		if err != nil {
			return fmt.Errorf("chat: error insert chat users in membri : %w", err)
		}

	}
	return nil
}
func (db *appdbimpl) CheckNames(cs int, toFind string) ([]UtenteDb, error) {
	var utenti []UtenteDb
	_, err := db.Authentication(cs)
	if err != nil {
		return utenti, fmt.Errorf("error in authentication Checknames: %w", err)
	}

	// ----------------------------------

	rows, err := db.c.Query("SELECT * FROM utenti WHERE (username LIKE $1 || '%') AND (id!=$2)", toFind, cs)
	if err != nil {
		return utenti, fmt.Errorf("chat: error querying users: %w", err)
	}

	var cont int

	defer rows.Close()

	for rows.Next() {
		cont++
		var utente UtenteDb
		if err := rows.Scan(&utente.Id, &utente.Username, &utente.Propic); err != nil {
			return utenti, fmt.Errorf("chat: error scanning user: %w", err)
		}
		utenti = append(utenti, utente)
	}
	if cont == 0 {
		return utenti, fmt.Errorf("chat: no user found: %w", err)
	}
	if err := rows.Err(); err != nil {
		return utenti, fmt.Errorf("chat: error iterating over users: %w", err)
	}

	// ----------------------------------------------

	return utenti, nil
}

// func (db *appdbimpl) GetMyConversations(cs int) ([]ChatDb, error) {
// 	var chat []ChatDb
// 	_, err := db.Authentication(cs)
// 	if err != nil {
// 		return chat, fmt.Errorf("error in authentication GetConversation: %w", err)
// 	}
// 	// la query sotto ritorna gli id degli utenti con cui ha la chat l'utente connesso, che non siano gruppi e l'id della chat
// 	rows, err := db.c.Query("SELECT m.id_utenti, c.id FROM chat c JOIN membri m ON c.id=m.id_chat WHERE m.id_utenti!=$1 AND c.id IN(SELECT  c.id from chat c JOIN membri m ON c.id=m.id_chat WHERE m.id_utenti=$1 AND gruppo=0) AND gruppo=0;", cs)
// 	if err != nil {
// 		return chat, fmt.Errorf("chat: error querying users: %w", err)
// 	}
// 	// questa query ritorna tutti i dati della join tra membri e chat dove l'utente appartiene al gruppo
// 	rows2, err := db.c.Query("SELECT * from chat c JOIN membri m ON c.id=m.id_chat WHERE m.id_utenti=$1 AND gruppo=1;", cs)
// 	if err != nil {
// 		return chat, fmt.Errorf("chat: error querying users: %w", err)
// 	}

// }
