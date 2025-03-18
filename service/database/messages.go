package database

import "fmt"

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

func (db *appdbimpl) CommentMessage(cs int, id_mes int) error {

	return nil
}
