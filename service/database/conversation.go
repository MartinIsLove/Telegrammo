package database

import (
	"database/sql"
	"fmt"
	"time"
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

func (db *appdbimpl) GetMyConversations(cs int) ([]ChatUtenteDb, error) {
	var chat []ChatUtenteDb
	_, err := db.Authentication(cs)
	if err != nil {
		return []ChatUtenteDb{}, fmt.Errorf("error in authentication GetConversation: %w", err)
	}
	// la query sotto ritorna gli id degli utenti con cui ha la chat l'utente connesso, che non siano gruppi e l'id della chat
	rows, err := db.c.Query("SELECT u.username, u.propic, m.id_utenti, c.id FROM chat c JOIN membri m ON c.id=m.id_chat JOIN utenti u ON u.id = m.id_utenti WHERE m.id_utenti!=$1 AND c.id IN(SELECT  c.id from chat c JOIN membri m ON c.id=m.id_chat WHERE m.id_utenti=$1 AND gruppo=0) AND gruppo=0;", cs)
	if err != nil {
		return []ChatUtenteDb{}, fmt.Errorf("chat: error querying users: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var c ChatUtenteDb
		if err := rows.Scan(&c.Nome, &c.Propic, &c.Id, &c.IdChat); err != nil {
			return []ChatUtenteDb{}, fmt.Errorf("chat: error scanning user: %w", err)
		}
		chat = append(chat, c)
	}
	if err := rows.Err(); err != nil {
		return []ChatUtenteDb{}, fmt.Errorf("chat: error iterating over users: %w", err)
	}

	// questa query ritorna tutti i dati della join tra membri e chat dove l'utente appartiene al gruppo
	rows2, err := db.c.Query("SELECT  c.id AS id, c.propic, c.nome, c.gruppo from chat c JOIN membri m ON c.id=m.id_chat JOIN utenti u ON u.id=m.id_utenti WHERE m.id_utenti=$1 AND gruppo=1;", cs)
	if err != nil {
		return []ChatUtenteDb{}, fmt.Errorf("chat: error querying users: %w", err)
	}

	defer rows2.Close()

	for rows2.Next() {
		var c ChatUtenteDb
		if err := rows2.Scan(&c.IdChat, &c.Propic, &c.Nome, &c.Gruppo); err != nil {
			return []ChatUtenteDb{}, fmt.Errorf("chat: error scanning user: %w", err)
		}
		chat = append(chat, c)
	}
	if err := rows2.Err(); err != nil {
		return []ChatUtenteDb{}, fmt.Errorf("chat: error iterating over users: %w", err)
	}
	fmt.Println("funz")

	for i := range chat {
		var lastMsg, username string
		var id int
		var tmp time.Time
		c := &chat[i]
		rows3 := db.c.QueryRow("SELECT m.testo, m.mittente, u.username, m.data FROM chat c JOIN messaggi_di_chat mdc ON c.id=mdc.id_chat JOIN messaggi m ON mdc.id_messaggio=m.id JOIN membri me ON me.id_chat=c.id JOIN utenti u ON u.id=me.id_utenti WHERE c.id=$1 ORDER BY m.data DESC LIMIT 1;", c.IdChat)
		err := rows3.Scan(&lastMsg, &id, &username, &tmp)
		if err == sql.ErrNoRows {

			c.LastMSg = ""
			c.Data = time.Time{}

		} else if err != nil {
			return []ChatUtenteDb{}, fmt.Errorf("chat: error querying users: %w", err)

		} else {
			if len(lastMsg) > 100 {
				c.LastMSg = lastMsg[:100] + "..."
			} else {
				c.LastMSg = lastMsg
			}
			c.Id = id
			c.Username = username
			c.Data = tmp
		}
	}
	// for _, c := range chat {
	// 	fmt.Printf("il nome della chat e': %s Username di chi ha inviato l'ultimo messaggio: %s, Propic: %s, UserId di chi ha inviato l'ultimo messaggio: %d, ChatId: %d , lastmsg: %s, data: %s, e' un gruppo: %t \n", c.Nome, c.Username, c.Propic, c.Id, c.IdChat, c.LastMSg, c.Data.Format(time.RFC3339), c.Gruppo)
	// }

	return chat, nil
}

func (db *appdbimpl) GetConversation(cs int, id_chat int) ([]MessDb, error) {
	var mess []MessDb
	_, err := db.Authentication(cs)
	if err != nil {
		return []MessDb{}, fmt.Errorf("error in authentication GetConversation: %w", err)
	}
	rows, err := db.c.Query("SELECT m.testo, m.mittente, u.username, m.data, m.image, m.id FROM messaggi m JOIN messaggi_di_chat d ON d.id_messaggio=m.id JOIN chat c ON c.id=d.id_chat JOIN utenti u ON u.id=m.mittente WHERE c.id=$1;", id_chat)
	if err != nil {
		return []MessDb{}, fmt.Errorf("chat: error querying users: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var c MessDb
		if err := rows.Scan(&c.Testo, &c.IdMitt, &c.Nome, &c.Data, &c.Photo, &c.IdMess); err != nil {
			return []MessDb{}, fmt.Errorf("chat: error scanning user: %w", err)
		}
		mess = append(mess, c)
	}
	if err := rows.Err(); err != nil {
		return []MessDb{}, fmt.Errorf("chat: error iterating over users: %w", err)
	}
	return mess, nil
}

func (db *appdbimpl) SendMessage(cs int, id_chat int, message string, photo []byte) error {
	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("error in authentication GetConversation: %w", err)
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
