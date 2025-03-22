package database

import (
	"fmt"
)

func (db *appdbimpl) SendMessage(cs int, id_chat int, message string, photo []byte) error {
	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("error in authentication SendMessage: %w", err)
	}
	str_mess, err := db.c.Exec("INSERT INTO messaggi (image, mittente, testo) VALUES ($1, $2, $3);", photo, cs, message)
	if err != nil {
		return fmt.Errorf("chat: error insert chat in table chat: %w", err)
	}

	id_mess, err := str_mess.LastInsertId()
	if err != nil {
		return fmt.Errorf("chat: error catch number of rows from query: %w", err)
	}

	_, err = db.c.Exec("INSERT INTO messaggi_di_chat (id_chat, id_messaggio) VALUES  ($1, $2)", id_chat, id_mess)
	if err != nil {
		return fmt.Errorf("chat: error insert messaggio in messagg_di_chat : %w", err)
	}

	return nil
}

func (db *appdbimpl) CommentMessage(cs int, id_mes int, emoji string, id_chat int) error {

	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("error in authentication SendMessage: %w", err)
	}

	var num_righe int64

	err = db.c.QueryRow("SELECT count(id_utenti) FROM membri WHERE id_utenti=$1 AND id_chat=$2", cs, id_chat).Scan(&num_righe)

	if err != nil {
		return fmt.Errorf("CommentMessage: query error: %w", err)
	}
	if num_righe == 0 {
		return fmt.Errorf("CommentMessage: you don't are in the chat %w", err)
	}

	err = db.c.QueryRow("SELECT count(id_utente) FROM emoticon WHERE id_utente=$1 AND id_messaggio=$2", cs, id_mes).Scan(&num_righe)

	if err != nil {
		return fmt.Errorf("CommentMessage: query error: %w", err)
	}

	if num_righe == 1 {
		_, err = db.c.Exec("UPDATE emoticon SET emoji=$1 WHERE id_messaggio=$2", emoji, id_mes)
		if err != nil {
			return fmt.Errorf("CommentMessage error: database UPDATE not successful: %w", err)
		}
		return nil
	}

	_, err = db.c.Exec("INSERT INTO emoticon (id_utente, id_messaggio, emoji) VALUES  ($1, $2, $3)", cs, id_mes, emoji)
	if err != nil {
		return fmt.Errorf("CommentMessage: error insert emoji in table emoji : %w", err)
	}

	return nil
}
