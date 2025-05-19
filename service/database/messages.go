package database

import (
	"fmt"
)

// inserisce un messaggio nel database
func (db *appdbimpl) SendMessage(cs int, id_chat int, message string, photo []byte, id_reply int) error {
	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("error in authentication SendMessage: %w", err)
	}

	var righe int

	// controlla se lo user che sta cercando di inviare il messaggio appartiene alla chat
	err = db.c.QueryRow("select count(u.id) from utenti u join membri m on m.id_utenti where id_utenti = $1 and id_chat = $2", cs, id_chat).Scan(&righe)
	if err != nil {
		return fmt.Errorf("SendMessage: error querying database: %w", err)
	}

	if righe == 0 {
		return fmt.Errorf("SendMessage: you don't belong to this chat")
	}

	// inserisce il messaggio nella tabella messaggi
	str_mess, err := db.c.Exec("INSERT INTO messaggi (image, mittente, testo) VALUES ($1, $2, $3);", photo, cs, message)
	if err != nil {
		return fmt.Errorf("SendMessage: error insert messagge in table messaggi: %w", err)
	}

	id_mess, err := str_mess.LastInsertId()
	if err != nil {
		return fmt.Errorf("SendMessage: error catch number of rows from query: %w", err)
	}

	// inserisce il messaggio e la chat nella tabella messaggi_di_chat
	_, err = db.c.Exec("INSERT INTO messaggi_di_chat (id_chat, id_messaggio, id_reply) VALUES  ($1, $2, $3)", id_chat, id_mess, id_reply)
	if err != nil {
		return fmt.Errorf("SendMessage: error insert messaggio in table messaggi_di_chat : %w", err)
	}

	return nil
}

// permette di fare il forward di un messaggio
func (db *appdbimpl) ForwardMessage(cs int, id_chat []int, id_mes int, id_utente int, id_utenti []int) error {
	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("error in authentication ForwardMessage: %w", err)
	}

	var tmp2 int

	// seleziono quanti utenti hanno nella tabella membri joinata con messaggi di chat hanno l'id dell'utente loggato e l'id del messaggio da forwardare
	err = db.c.QueryRow("SELECT COUNT(m.id_utenti) FROM chat c JOIN messaggi_di_chat mdc ON c.id=mdc.id_chat JOIN membri m ON m.id_chat=c.id WHERE mdc.id_messaggio=$1 AND m.id_utenti=$2", id_mes, cs).Scan(&tmp2)

	if err != nil {
		return fmt.Errorf("ForwardMessage: error select id_utenti from chat joined on messaggi di chat: %w", err)
	}

	// se non esiste allora non hai accesso al messaggio
	if tmp2 == 0 {
		return fmt.Errorf("ForwardMessage: error you don't belong to the chat")
	}

	var num_righe int
	var tmp int

	// per ogni chat in cui forwardare il messaggio
	for _, c := range id_chat {

		// verifico l'appartenenza alle chat in cui voglio forwardare il messaggio
		err = db.c.QueryRow("SELECT count(id_utenti) FROM membri WHERE id_utenti=$1 AND id_chat=$2", cs, c).Scan(&num_righe)

		if err != nil {
			return fmt.Errorf("ForwardMessage:2 error select id_utenti from membri: %w", err)
		}
		tmp += num_righe
	}

	if tmp != len(id_chat) {
		return fmt.Errorf("ForwardMessage: you don't belong to the group you want forward the message")
	}

	// per ogni chat in cui forwwardare
	for _, c := range id_chat {

		// inserisce il messaggio forwardato
		_, err = db.c.Exec("INSERT INTO messaggi_di_chat (id_chat, id_messaggio, id_forward, id_forw_mit,id_reply) VALUES ($1, $2, $3, $4, $5);", c, id_mes, id_utente, cs, -1)
		if err != nil {
			return fmt.Errorf("ForwardMessage: error insert message in table messaggi_di_chat: %w", err)
		}
	}

	for _, c := range id_utenti {

		err = db.c.QueryRow(`
			SELECT COUNT(c.id)
			FROM chat c
			JOIN membri m1 ON m1.id_chat = c.id AND m1.id_utenti = $1
			JOIN membri m2 ON m2.id_chat = c.id AND m2.id_utenti = $2
			WHERE c.gruppo = 0
		`, cs, c).Scan(&num_righe)
		if err != nil {
			return fmt.Errorf("ForwardMessage: error querying chat count: %w", err)
		}
		if num_righe == 0 {
			res, err := db.c.Exec("INSERT INTO chat (nome, propic , gruppo) VALUES (null, null, 0);")
			if err != nil {
				return fmt.Errorf("ForwardMessage: error insert chat in table chat: %w", err)
			}
			chatID, err := res.LastInsertId()
			if err != nil {
				return fmt.Errorf("ForwardMessage: error getting last insert id: %w", err)
			}
			_, err = db.c.Exec("INSERT INTO membri (id_utenti,id_chat) VALUES ($1,$2),($3,$2);", cs, chatID, c)
			if err != nil {
				return fmt.Errorf("ForwardMessage: error insert chat in table chat: %w", err)
			}

			_, err = db.c.Exec("INSERT INTO messaggi_di_chat (id_chat, id_messaggio, id_forward, id_forw_mit,id_reply) VALUES ($1, $2, $3, $4, $5);", chatID, id_mes, id_utente, cs, -1)
			if err != nil {
				return fmt.Errorf("ForwardMessage: error insert message in table messaggi_di_chat: %w", err)
			}
		} else {
			err = db.c.QueryRow(`SELECT c.id
			FROM chat c
			JOIN membri m1 ON m1.id_chat = c.id AND m1.id_utenti = $1
			JOIN membri m2 ON m2.id_chat = c.id AND m2.id_utenti = $2
			WHERE c.gruppo = 0
			LIMIT 1;`, cs, c).Scan(&num_righe)

			if err != nil {
				return fmt.Errorf("ForwardMessage:1 error select id_utenti from membri: %w", err)
			}

			_, err = db.c.Exec("INSERT INTO messaggi_di_chat (id_chat, id_messaggio, id_forward, id_forw_mit,id_reply) VALUES ($1, $2, $3, $4, $5);", num_righe, id_mes, id_utente, cs, -1)
			if err != nil {
				return fmt.Errorf("ForwardMessage: error insert message in table messaggi_di_chat: %w", err)
			}
		}
	}

	return nil
}

// permette di mettere un emoji ad un messaggio
func (db *appdbimpl) CommentMessage(cs int, id_mes int, emoji string, id_chat int) error {

	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("error in authentication CommentMessage: %w", err)
	}

	var num_righe int64

	// vedo se l'utente fa parte della chat per poter commentare il messaggio
	err = db.c.QueryRow("SELECT count(id_utenti) FROM membri WHERE id_utenti=$1 AND id_chat=$2", cs, id_chat).Scan(&num_righe)

	if err != nil {
		return fmt.Errorf("CommentMessage: error select id_utenti form membri: %w", err)
	}
	if num_righe == 0 {
		return fmt.Errorf("CommentMessage: you don't belong to the chat %w", err)
	}

	// per verificare se ho gia inserito un commento per quel messaggio
	err = db.c.QueryRow("SELECT count(id_utente) FROM emoticon WHERE id_utente=$1 AND id_messaggio=$2", cs, id_mes).Scan(&num_righe)

	if err != nil {
		return fmt.Errorf("CommentMessage: error select id_utente from emoticon: %w", err)
	}

	// in caso positivo lo modifico
	if num_righe == 1 {
		_, err = db.c.Exec("UPDATE emoticon SET emoji=$1 WHERE id_messaggio=$2 AND id_utente=$3", emoji, id_mes, cs)
		if err != nil {
			return fmt.Errorf("CommentMessage: error database UPDATE not successful: %w", err)
		}
		return nil
	}

	// altrimenti lo aggiungo
	_, err = db.c.Exec("INSERT INTO emoticon (id_utente, id_messaggio, emoji) VALUES  ($1, $2, $3)", cs, id_mes, emoji)
	if err != nil {
		return fmt.Errorf("CommentMessage: error insert emoji in table emoticon : %w", err)
	}

	return nil
}

// elimina il commento ad un messaggio
func (db *appdbimpl) UncommentMessage(cs int, id_mes int, id_chat int) error {

	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("error in authentication UncommentMessage: %w", err)
	}

	var num_righe int64

	// verifica la presenza nella chat nella quale e' presente il messaggio da uncommentare
	err = db.c.QueryRow("SELECT count(id_utenti) FROM membri WHERE id_utenti=$1 AND id_chat=$2", cs, id_chat).Scan(&num_righe)

	if err != nil {
		return fmt.Errorf("UncommentMessage: error select id_utenti from membri: %w", err)
	}
	if num_righe == 0 {
		return fmt.Errorf("UncommentMessage: you don't belong to the chat %w", err)
	}

	// conta nella tabella emoticon quanti id_utenti e id_messaggio ci sono uguali a quelli che ho inserito
	err = db.c.QueryRow("SELECT count(id_utente) FROM emoticon WHERE id_utente=$1 AND id_messaggio=$2", cs, id_mes).Scan(&num_righe)

	if err != nil {
		return fmt.Errorf("UncommentMessage: error select id_utente from emoticon: %w", err)
	}

	// se e' stato trovato puo' essere eliminato
	if num_righe == 1 {
		// elimina il commento
		_, err = db.c.Exec("DELETE FROM emoticon WHERE id_utente=$1 AND id_messaggio=$2", cs, id_mes)
		if err != nil {
			return fmt.Errorf("UncommentMessage: error database DELETE not successful: %w", err)
		}
		return nil
	}

	return nil

}

// cancella un messaggio
func (db *appdbimpl) DeleteMessage(cs int, idMes int, idForw int, idChat int) error {
	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("error in authentication DeleteMessage: %w", err)
	}

	var is_forw int

	// conto quanti messaggi nella tabella messaggi_di_chat hanno lo stesso id_messaggio del messaggio in questione, per capire se e' stato o no forwardato
	err = db.c.QueryRow("SELECT COUNT(id) FROM messaggi_di_chat WHERE id_messaggio=$2", idMes).Scan(&is_forw)
	if err != nil {
		return fmt.Errorf("DeleteMessage:  error select id from messaggi_di_chat: %w", err)
	}

	// se e' 0 allora il messaggio non esiste
	if is_forw == 0 {
		return fmt.Errorf("DeleteMessage: error database DELETE not successful message don't find")
	}

	// se e' 1 allora puo' essere eliminato del tutto
	if is_forw == 1 {
		_, err = db.c.Exec("DELETE FROM messaggi WHERE id=$1", idMes)
		if err != nil {
			return fmt.Errorf("DeleteMessage: error database DELETE not successful: %w", err)
		}
		_, err = db.c.Exec("DELETE FROM messaggi_di_chat WHERE id=$1", idForw)
		if err != nil {
			return fmt.Errorf("DeleteMessage: error database DELETE not successful: %w", err)
		}
	}

	// se e' maggiore di 1 significa che e' stato forwardato e che quindi non si puo eliminare l'istanza principale nella tabella messaggi
	if is_forw > 1 {
		_, err = db.c.Exec("DELETE FROM messaggi_di_chat WHERE id=$1", idForw)
		if err != nil {
			return fmt.Errorf("DeleteMessage: error database DELETE not successful: %w", err)
		}
	}

	// per ricordare questa query controlla se il messaggio e' stato forwardato o no
	// nel caso in cui non lo e' stato basta cancellare il messaggio principale
	// nel caso in cui e' stato forwardato bisogna cancellare l'istanza del messaggio
	return nil
}
