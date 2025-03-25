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
func (db *appdbimpl) ForwardMessage(cs int, id_chat []int, id_mes int) error {
	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("error in authentication SendMessage: %w", err)
	}
	var tmp2 int
	err = db.c.QueryRow("SELECT COUNT(m.id_utenti) FROM chat c JOIN messaggi_di_chat mdc ON c.id=mdc.id_chat JOIN membri m ON m.id_chat=c.id WHERE mdc.id_messaggio=$1 AND m.id_utenti=$2", cs, id_mes).Scan(&tmp2)

	if err != nil {
		return fmt.Errorf("CommentMessage: query error: %w", err)
	}

	if tmp2 != 1 {
		return fmt.Errorf("non hai accesso al messaggio da forwardare")
	}

	var num_righe int
	var tmp int
	for _, c := range id_chat {
		err = db.c.QueryRow("SELECT count(id_utenti) FROM membri WHERE id_utenti=$1 AND id_chat=$2", cs, c).Scan(&num_righe)

		if err != nil {
			return fmt.Errorf("CommentMessage: query error: %w", err)
		}
		tmp = tmp + num_righe
	}

	if tmp != len(id_chat) {
		return fmt.Errorf("ForwardMessage: non appartieni a una o piu chat di quelle che hai selezionato")
	}

	for _, c := range id_chat {
		_, err = db.c.Exec("INSERT INTO messaggi_di_chat (id_chat, id_messaggio) VALUES ($1, $2);", c, id_mes)
		if err != nil {
			return fmt.Errorf("ForwardMessage: error insert chat in table chat: %w", err)
		}
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
		_, err = db.c.Exec("UPDATE emoticon SET emoji=$1 WHERE id_messaggio=$2 AND id_utente=$3", emoji, id_mes, cs)
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
func (db *appdbimpl) UncommentMessage(cs int, id_mes int, id_chat int) error {

	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("error in authentication SendMessage: %w", err)
	}

	var num_righe int64

	err = db.c.QueryRow("SELECT count(id_utenti) FROM membri WHERE id_utenti=$1 AND id_chat=$2", cs, id_chat).Scan(&num_righe)

	if err != nil {
		return fmt.Errorf("UncommentMessage: query error: %w", err)
	}
	if num_righe == 0 {
		return fmt.Errorf("UncommentMessage: you don't are in the chat %w", err)
	}
	err = db.c.QueryRow("SELECT count(id_utente) FROM emoticon WHERE id_utente=$1 AND id_messaggio=$2", cs, id_mes).Scan(&num_righe)

	if err != nil {
		return fmt.Errorf("UncommentMessage: query error: %w", err)
	}
	if num_righe == 1 {
		_, err = db.c.Exec("DELETE FROM emoticon WHERE id_utente=$1 AND id_messaggio=$2", cs, id_mes)
		if err != nil {
			return fmt.Errorf("UncommentMessage error: database DELETE not successful: %w", err)
		}
		return nil
	}

	return nil

}
