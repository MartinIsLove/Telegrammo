package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// ritorna 1 se una chat(non gruppo) gia esiste, 0 altrimenti
func (db *appdbimpl) IsChatDuplicated(cs int, id int) (bool, error) {
	var num_righe int64
	// seleziona la chat(non gruppo) con partecipanti cs e id dati nei parametri della funzione
	err := db.c.QueryRow("SELECT COUNT(chat.id) AS righe FROM chat JOIN membri ON chat.id=membri.id_chat WHERE chat.gruppo=FALSE AND membri.id_utenti IN ($1 , $2) GROUP BY chat.id HAVING COUNT (DISTINCT membri.id_utenti)=2", cs, id).Scan(&num_righe)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return true, fmt.Errorf("IsChatDuplicated: error checking chat duplicated in database: %w", err)
	}

	if num_righe == 0 {
		return false, nil
	}
	return true, fmt.Errorf("IsChatDuplicated: error this chat already exist: %w", err)
}

// crea una chat(non gruppo)
func (db *appdbimpl) CreateChat(cs int, id int) (int, error) {

	_, err := db.Authentication(cs)
	if err != nil {
		return -1, fmt.Errorf("error in authentication CreateChat: %w", err)
	}

	// verifica se la chat e' duplicata
	Isduplicated, err := db.IsChatDuplicated(cs, id)

	if err != nil {
		return -1, err
	}

	// se non lo e'
	if !Isduplicated {
		// imposta il valore gruppo come false (essendo una chat)
		str_chat, err := db.c.Exec("INSERT INTO chat (gruppo) VALUES (FALSE)")
		if err != nil {
			return -1, fmt.Errorf("CreateChat: error insert chat in table chat: %w", err)
		}

		// prendo l'id della chat appena creata
		id_chat, err := str_chat.LastInsertId()
		if err != nil {
			return -1, fmt.Errorf("CreateChat: error catch number of rows from query: %w", err)
		}

		// inserisco nell'id della chat appena creata nella tabella membri i due id degli utenti che ne fanno parte
		_, err = db.c.Exec("INSERT INTO membri (id_chat, id_utenti) VALUES  ($1, $2), ($1, $3)", id_chat, cs, id)
		if err != nil {
			return -1, fmt.Errorf("CreateChat: error insert chat users in membri : %w", err)
		}
		return int(id_chat), nil
	} else {
		return -1, err
	}

}

// ritorna i match tra i caratteri inseriti e i primi caratteri dei nomi delle chat
func (db *appdbimpl) CheckChatNames(cs int, toFind string) ([]ChatUtenteDb, error) {

	var chat []ChatUtenteDb
	_, err := db.Authentication(cs)
	if err != nil {
		return []ChatUtenteDb{}, fmt.Errorf("error in authentication CheckChatNames: %w", err)
	}

	// la query sotto ritorna gli id degli utenti con cui ha la chat l'utente connesso, che non siano gruppi e l'id della chat
	rows, err1 := db.c.Query("SELECT u.username, u.propic, m.id_utenti, c.id FROM chat c JOIN membri m ON c.id=m.id_chat JOIN utenti u ON u.id = m.id_utenti WHERE (u.username LIKE $1 || '%') AND m.id_utenti!=$2 AND c.id IN(SELECT  c.id from chat c JOIN membri m ON c.id=m.id_chat WHERE m.id_utenti=$2 AND gruppo=0) AND gruppo=0;", toFind, cs)
	if err1 != nil {
		return []ChatUtenteDb{}, fmt.Errorf("CheckChatNames: error querying users: %w", err)
	}

	defer rows.Close()

	hasRows := false

	for rows.Next() {
		hasRows = true
		var c ChatUtenteDb
		if err := rows.Scan(&c.Nome, &c.Propic, &c.Id, &c.IdChat); err != nil {
			return []ChatUtenteDb{}, fmt.Errorf(" CheckChatNames: error scanning user: %w", err)
		}

		chat = append(chat, c)
	}

	if err := rows.Err(); err != nil {
		return []ChatUtenteDb{}, fmt.Errorf("CheckChatNames: error iterating over users: %w", err)
	}

	// questa query ritorna tutti i dati della join tra membri e chat dove l'utente appartiene al gruppo
	rows2, err2 := db.c.Query("SELECT  c.id AS id, c.propic, c.nome, c.gruppo from chat c JOIN membri m ON c.id=m.id_chat JOIN utenti u ON u.id=m.id_utenti WHERE (c.nome LIKE $1 || '%') AND m.id_utenti=$2 AND gruppo=1;", toFind, cs)
	if err2 != nil {
		return []ChatUtenteDb{}, fmt.Errorf("CheckChatNames: error querying users: %w", err)
	}

	defer rows2.Close()

	hasRows2 := false

	for rows2.Next() {
		hasRows2 = true

		var c ChatUtenteDb
		if err := rows2.Scan(&c.IdChat, &c.Propic, &c.Nome, &c.Gruppo); err != nil {
			return []ChatUtenteDb{}, fmt.Errorf("CheckChatNames: error scanning user: %w", err)
		}
		chat = append(chat, c)
	}
	if err := rows2.Err(); err != nil {
		return []ChatUtenteDb{}, fmt.Errorf("CheckChatNames: error iterating over users: %w", err)
	}

	if !hasRows && !hasRows2 {
		return []ChatUtenteDb{}, fmt.Errorf("CheckChatNames: no chats or groups find: %w", err)

	}

	return chat, nil
}

// ritorna tutte le conversazioni dell'utente loggato
func (db *appdbimpl) GetMyConversations(cs int) ([]ChatUtenteDb, error) {
	var chat []ChatUtenteDb
	_, err := db.Authentication(cs)
	if err != nil {
		return []ChatUtenteDb{}, fmt.Errorf("error in authentication GetConversations: %w", err)
	}
	// la query sotto ritorna gli id degli utenti con cui ha la chat l'utente connesso, che non siano gruppi e l'id della chat
	rows, err1 := db.c.Query("SELECT u.username, u.propic, m.id_utenti, c.id FROM chat c JOIN membri m ON c.id=m.id_chat JOIN utenti u ON u.id = m.id_utenti WHERE m.id_utenti!=$1 AND c.id IN(SELECT  c.id from chat c JOIN membri m ON c.id=m.id_chat WHERE m.id_utenti=$1 AND gruppo=0) AND gruppo=0;", cs)
	if err1 != nil && !errors.Is(err1, sql.ErrNoRows) {
		return []ChatUtenteDb{}, fmt.Errorf("GetMyConversations: error querying chats: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var c ChatUtenteDb
		if err := rows.Scan(&c.Nome, &c.Propic, &c.Id, &c.IdChat); err != nil {
			return []ChatUtenteDb{}, fmt.Errorf(" GetMyconversations: error scanning chats: %w", err)
		}
		chat = append(chat, c)
	}
	if err := rows.Err(); err != nil {
		return []ChatUtenteDb{}, fmt.Errorf("GetMyConversations: error iterating over chats: %w", err)
	}

	// questa query ritorna tutti i dati della join tra membri e chat dove l'utente appartiene al gruppo
	rows2, err2 := db.c.Query("SELECT  c.id AS id, c.propic, c.nome, c.gruppo from chat c JOIN membri m ON c.id=m.id_chat JOIN utenti u ON u.id=m.id_utenti WHERE m.id_utenti=$1 AND gruppo=1;", cs)
	if err2 != nil && !errors.Is(err2, sql.ErrNoRows) {
		return []ChatUtenteDb{}, fmt.Errorf("GetMyConversations: error querying chats: %w", err)
	}

	defer rows2.Close()

	for rows2.Next() {
		var c ChatUtenteDb
		if err := rows2.Scan(&c.IdChat, &c.Propic, &c.Nome, &c.Gruppo); err != nil {
			return []ChatUtenteDb{}, fmt.Errorf("GetMyConversations: error scanning chats: %w", err)
		}
		chat = append(chat, c)
	}
	if err := rows2.Err(); err != nil {
		return []ChatUtenteDb{}, fmt.Errorf("GetMyConversations: error iterating over chats: %w", err)
	}

	if errors.Is(err1, sql.ErrNoRows) && errors.Is(err2, sql.ErrNoRows) {
		return []ChatUtenteDb{}, fmt.Errorf("GetMyConversations: no chats or groups find: %w", err)

	}

	for i := range chat {
		var lastMsg, username string
		var id int
		var tmp time.Time
		c := &chat[i]
		// prendo gli ultimi messaggi inviati a ogni chat presa precedentemente
		rows3 := db.c.QueryRow("SELECT m.testo, m.mittente, m.data FROM chat c JOIN messaggi_di_chat mdc ON c.id=mdc.id_chat JOIN messaggi m ON mdc.id_messaggio=m.id JOIN membri me ON me.id_chat=c.id JOIN utenti u ON u.id=me.id_utenti WHERE c.id=$1 ORDER BY m.data DESC LIMIT 1;", c.IdChat)
		err := rows3.Scan(&lastMsg, &id, &tmp)
		if errors.Is(err, sql.ErrNoRows) {

			c.LastMSg = ""
			c.Data = time.Time{}

		} else if err != nil {
			return []ChatUtenteDb{}, fmt.Errorf("GetMyConversations: error querying chats: %w", err)

		} else {
			if len(lastMsg) > 100 {
				c.LastMSg = lastMsg[:100] + "..."
			} else {
				c.LastMSg = lastMsg
			}
			c.Id = id
			c.Data = tmp
		}
		// seleziono l'username di chi ha inviato l'ultimo messaggio alla chat
		rows4 := db.c.QueryRow("SELECT username FROM utenti WHERE id=$1 LIMIT 1", id)
		err = rows4.Scan(&username)
		if errors.Is(err, sql.ErrNoRows) {

		} else if err != nil {
			return []ChatUtenteDb{}, fmt.Errorf("GetMyConversations: error querying users: %w", err)

		}
		c.Username = username
	}

	return chat, nil
}

// ritorna le informazioni della chat id_chat
func (db *appdbimpl) GetConversation(cs int, id_chat int) (bool, string, []MessDb, error) {
	var mess []MessDb
	_, err := db.Authentication(cs)
	if err != nil {
		return false, "", []MessDb{}, fmt.Errorf("error in authentication GetConversation: %w", err)
	}

	var belong int

	// controllo l'appartenenza al gruppo
	err = db.c.QueryRow("SELECT COUNT(id_utenti) FROM membri WHERE id_utenti=$1 AND id_CHAT=$2", cs, id_chat).Scan(&belong)
	if err != nil {
		return false, "", []MessDb{}, fmt.Errorf("GetConversation: error querying database: %w", err)
	}
	if belong == 0 {
		return false, "", []MessDb{}, fmt.Errorf("GetConversation: you don't belong to this group %w", err)
	}

	// inserisce nella tabella accessi_chat l'accesso avvenuto a quella chat e nel caso in cui gia esista l'accesso di quell'utente a quella chat lo sovrascrive
	_, err = db.c.Exec("INSERT INTO accessi_chat (id_utente, id_chat, data) VALUES ($1, $2, CURRENT_TIMESTAMP) ON CONFLICT (id_utente, id_chat) DO UPDATE SET data = CURRENT_TIMESTAMP WHERE accessi_chat.id_utente = $1 AND accessi_chat.id_chat = $2", cs, id_chat)
	if err != nil {
		return false, "", []MessDb{}, fmt.Errorf("GetConversation: error update accessi_chat: %w", err)
	}

	var group bool

	// verifico se la chat e' un gruppo oppure no
	err = db.c.QueryRow("SELECT gruppo FROM chat WHERE id=$1 LIMIT 1;", id_chat).Scan(&group)

	if err != nil {
		return false, "", []MessDb{}, fmt.Errorf("GetConversation: error selecting if group: %w", err)
	}
	var nomeChat string

	// se lo e'
	if group {
		// prendo il nome del gruppo
		err = db.c.QueryRow("SELECT nome FROM chat WHERE id=$1 LIMIT 1;", id_chat).Scan(&nomeChat)
		if err != nil {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: error selecting group name: %w", err)
		}
	} else {
		// prendo il nome della chat(il nome dell'altro utente che fa parte della chat)
		err = db.c.QueryRow("SELECT utenti.username FROM chat c JOIN membri ON c.id = membri.id_chat JOIN utenti ON membri.id_utenti = utenti.id WHERE c.id=$1 AND utenti.id != $2 LIMIT 1;", id_chat, cs).Scan(&nomeChat)
		if err != nil {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: error selecting group name: %w", err)
		}
	}
	// seleziono tutto i messaggi della chat, compresi se fa parte di un gruppo, testo, id mittente, username mittente, data dell'invio, l'immagine(se presente), l'id del messaggio e l'id del messaggio di chat
	rows, err := db.c.Query("SELECT c.gruppo, m.testo, m.mittente, u.username, m.data, m.image, m.id,d.id FROM messaggi m JOIN messaggi_di_chat d ON d.id_messaggio=m.id JOIN chat c ON c.id=d.id_chat JOIN utenti u ON u.id=m.mittente WHERE c.id=$1 ORDER BY m.data;", id_chat)
	if errors.Is(err, sql.ErrNoRows) {
		return false, "", []MessDb{}, fmt.Errorf("GetConversation: no messages found: %w", err)
	}
	if err != nil {
		return false, "", []MessDb{}, fmt.Errorf("GetConversation: error querying messages: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var c MessDb
		if err := rows.Scan(&c.Gruppo, &c.Testo, &c.IdMitt, &c.Nome, &c.Data, &c.Photo, &c.IdMess, &c.IdForward); err != nil {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: error scanning messages: %w", err)
		}

		var forwarderUsername string
		var forwarderId int
		var forwardDate time.Time
		var forwardIdMit int
		var forwarderUsernameMit string
		var cont int = 0

		// seleziono i messaggi forwardati tra i messaggi presi prima abbiamo l'id del messaggio di chat uguale e l'id del messaggio uguale, prendendo in questo modo le informazioni del forward per quel dato messaggio
		row2, err := db.c.Query("SELECT u.username, u.id, m.forward_date, m.id_forw_mit, t.username FROM utenti u JOIN messaggi_di_chat m ON m.id_forward = u.id JOIN utenti t ON m.id_forw_mit=t.id WHERE m.id_messaggio = $1 AND m.id_forward NOT NULL AND m.id=$2", c.IdMess, c.IdForward)
		if errors.Is(err, sql.ErrNoRows) {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: no messages found: %w", err)
		}
		if err != nil {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: error querying messages: %w", err)
		}

		defer row2.Close()

		for row2.Next() {
			if err := row2.Scan(&forwarderUsername, &forwarderId, &forwardDate, &forwardIdMit, &forwarderUsernameMit); err != nil {
				return false, "", []MessDb{}, fmt.Errorf("GetConversation: error scanning messages: %w", err)
			}

			c.ForwardUsername = forwarderUsername
			c.ForwardId = forwarderId
			c.ForwardDate = forwardDate
			c.ForwardIdMit = forwardIdMit
			c.ForwardUsernameMit = forwarderUsernameMit
			mess = append(mess, c)
			cont++
		}

		if err := row2.Err(); err != nil {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: error iterating messages: %w", err)
		}

		if cont == 0 {
			c.ForwardUsername = forwarderUsername
			c.ForwardId = -1
			c.ForwardDate = forwardDate
			c.ForwardIdMit = -1
			c.ForwardUsernameMit = forwarderUsernameMit
			mess = append(mess, c)
		}
	}
	if err := rows.Err(); err != nil {
		return false, "", []MessDb{}, fmt.Errorf("GetConversation: error iterating over messages: %w", err)
	}
	// avendo aggiunto i messaggi forwardati ora vanno ordinati per data di forward rispetto alla data degli altri messaggi
	sort.SliceStable(mess, func(i, j int) bool {
		if mess[i].ForwardId != -1 && mess[j].ForwardId != -1 {
			return mess[i].ForwardDate.Before(mess[j].ForwardDate)
		} else if mess[i].ForwardId != -1 {
			return mess[i].ForwardDate.Before(mess[j].Data)
		} else if mess[j].ForwardId != -1 {
			return mess[i].Data.Before(mess[j].ForwardDate)
		} else {
			return mess[i].Data.Before(mess[j].Data)
		}
	})

	var result []MessDb
	var lastDate time.Time

	// per ogni messaggio (serve per verificare la visualizzazione del messaggio da parte di tutti gli utenti)
	for _, m := range mess {
		var emoji string
		var idUtenti []int

		// seleziono tutti gli id utenti della chat e li metto in un array
		row3, err := db.c.Query("SELECT id_utenti FROM membri WHERE id_utenti != $1 AND id_chat=$2", cs, id_chat)
		if errors.Is(err, sql.ErrNoRows) {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: no membri found: %w", err)
		}
		if err != nil {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: error querying users: %w", err)
		}

		defer row3.Close()

		for row3.Next() {
			var id int
			if err := row3.Scan(&id); err != nil {
				return false, "", []MessDb{}, fmt.Errorf("GetConversation: error scanning id_utenti: %w", err)
			}
			idUtenti = append(idUtenti, id)
		}

		if err := row3.Err(); err != nil {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: error iterating over id_utenti: %w", err)
		}

		var visualCount int
		// conto il numero di tutti gli id utenti diversi tra loro che fanno parte della chat, l'utente non e' il mittente e la data dell'accesso e' successiva all'invio del messaggio
		err = db.c.QueryRow(`SELECT COUNT(DISTINCT ac.id_utente) FROM accessi_chat ac JOIN membri mb ON mb.id_utenti = ac.id_utente WHERE ac.id_chat = $1 AND ac.id_utente != $2 AND ac.data > $3 AND mb.id_chat = ac.id_chat`, id_chat, m.IdMitt, m.Data).Scan(&visualCount)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: error counting message views: %w", err)
		}

		// se le visualizzazioni sono maggiori o uguali agli utenti allora il messaggio e' visualizzato, altrimenti no
		if visualCount >= len(idUtenti) {
			m.Visual = true
		} else {
			m.Visual = false
		}
		// -----------------------------------------------------

		// seleziono l'id del reply del messaggio in questione dalla tabella messaggi di chat, usando IdForward come id
		err = db.c.QueryRow("SELECT id_reply FROM messaggi_di_chat WHERE id=$1", m.IdForward).Scan(&m.Idreply)
		if err != nil {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: error querying messages: %w", err)
		}
		if m.Idreply > 0 {

			// prendo le informazioni del messaggio replyato
			err = db.c.QueryRow("SELECT m.testo, m.image, m.mittente, u.username FROM messaggi m JOIN utenti u ON u.id=m.mittente JOIN messaggi_di_chat md ON md.id_messaggio=m.id WHERE md.id=$1", m.Idreply).Scan(&m.TestoReply, &m.PhotoReply, &m.IdMitReply, &m.MitReply)

			// se non restituisce righe
			if errors.Is(err, sql.ErrNoRows) {

				// allora cancella il reply
				_, err := db.c.Exec("UPDATE messaggi_di_chat SET id_reply=$1 WHERE id=$2", -1, m.IdForward)
				if err != nil {
					return false, "", []MessDb{}, fmt.Errorf("GetConversation: error insert id_reply in table messaggi_di_chat: %w", err)
				}
				m.TestoReply = ""
				m.PhotoReply = nil
				m.IdMitReply = -1
				m.MitReply = ""
			} else if err != nil {
				return false, "", []MessDb{}, fmt.Errorf("GetConversation: error querying messages: %w", err)
			}

		}

		// seleziono l'emoji inviata da me a quel messaggio se esiste altrimenti la imposto vuota
		err = db.c.QueryRow("SELECT emoji FROM emoticon WHERE id_utente=$1 AND id_messaggio=$2", cs, m.IdMess).Scan(&emoji)
		if errors.Is(err, sql.ErrNoRows) {
			emoji = ""
		} else if err != nil {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: error querying emojis: %w", err)
		}

		var counts [5]int

		// conto per ogni tipo di emoji quante ne sono state lasciate per ciascun messaggio
		err = db.c.QueryRow("SELECT COALESCE(SUM(CASE WHEN emoji = '👠' AND id_messaggio=$1 THEN 1 ELSE 0 END),0) AS tacchi, COALESCE(SUM(CASE WHEN emoji = '❤️' AND id_messaggio=$1 THEN 1 ELSE 0 END),0) AS cuore, COALESCE(SUM(CASE WHEN emoji = '👍🏻' AND id_messaggio=$1 THEN 1 ELSE 0 END),0) AS dito_in_su, COALESCE(SUM(CASE WHEN emoji = '👌🏻' AND id_messaggio=$1 THEN 1 ELSE 0 END),0) AS ok, COALESCE(SUM(CASE WHEN emoji = '💅' AND id_messaggio=$1 THEN 1 ELSE 0 END),0) AS manicure FROM emoticon;", m.IdMess).Scan(&counts[0], &counts[1], &counts[2], &counts[3], &counts[4])
		if err != nil {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: error querying emojis: %w", err)
		}

		// aggiungo dei messaggi speciali, riconosciuti per la non presenza ne di testo ne di foto nei quali sono presenti solo le indicazioni di data, che servono per separare nella chat i diversi giorni
		localTime := m.Data.Local().Add(+1 * time.Hour)
		localMidnight := time.Date(localTime.Year(), localTime.Month(), localTime.Day(), 0, 0, 0, 0, localTime.Location())
		if lastDate.IsZero() || !localMidnight.Equal(lastDate) {
			dateMessage := MessDb{
				Data: localMidnight,
			}
			result = append(result, dateMessage)
			lastDate = localMidnight
		}
		if emoji != "" {
			m.MyEmoji = emoji
		}
		m.Emoji = counts[:]
		result = append(result, m)
	}

	return group, nomeChat, result, nil
}

// crea un gruppo
func (db *appdbimpl) CreateGroup(cs int, nome string, propic []byte, membri []int) (int, error) {
	_, err := db.Authentication(cs)
	if err != nil {
		return -1, fmt.Errorf("error in authentication CreateGroup: %w", err)
	}
	noPhotoPath := filepath.Join("/workspace/webui/src/assets/", "NoPhoto.png")
	noPhotoBytes, err := os.ReadFile(noPhotoPath)
	if err != nil {
		return -1, fmt.Errorf("CreateGroup: error reading noPhoto.png: %w", err)
	}
	if propic == nil {
		propic = noPhotoBytes
	}

	// inserisco il gruppo nella tabella chat
	str_group, err := db.c.Exec("INSERT INTO chat (nome, propic, gruppo) VALUES ($1, $2, TRUE)", nome, propic)
	if err != nil {
		return -1, fmt.Errorf("CreateGroup: error insert group in table chat: %w", err)
	}

	// recupero l'id del gruppo appena creato
	id_group, err := str_group.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("CreateGroup: error catch number of rows from query: %w", err)
	}

	// inserisco l'utente loggato nel gruppo
	_, err = db.c.Exec("INSERT INTO membri (id_utenti, id_chat) VALUES ($1, $2)", cs, id_group)
	if err != nil {
		return -1, fmt.Errorf("CreateGroup: error insert (user, chat) in table membri: %w", err)
	}
	for _, m := range membri {
		// controllo che lo user che sto inserendo esista
		value := db.UserExist(m)

		// se esiste
		if value {

			// inserisco lo user(lo fa per ogni user in membri) nella tabella membri
			_, err := db.c.Exec("INSERT INTO membri (id_utenti, id_chat) VALUES ($1, $2)", m, id_group)
			if err != nil {
				return -1, fmt.Errorf("CreateGroup: error insert group in table chat: %w", err)
			}
		}
	}

	return int(id_group), nil
}

// rimuove l'utente loggato da un gruppo
func (db *appdbimpl) LeaveGroup(cs int, idChat int) error {

	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("error in authentication LeaveGroup: %w", err)
	}

	var if_group int

	// verifico se la chat e' un gruppo
	err = db.c.QueryRow("SELECT COUNT(id_utenti) FROM membri m JOIN chat c ON c.id=m.id_chat WHERE m.id_utenti=$1 AND m.id_CHAT=$2 AND c.gruppo=1", cs, idChat).Scan(&if_group)
	if err != nil {
		return fmt.Errorf("LeaveGroup: error querying database: %w", err)
	}
	if if_group == 0 {
		return fmt.Errorf("LeaveGroup: this isn't a group: %w", err)
	}

	if_group = 0
	// controllo l'appartenenza al gruppo
	err = db.c.QueryRow("SELECT COUNT(id_utenti) FROM membri WHERE id_utenti=$1 AND id_CHAT=$2", cs, idChat).Scan(&if_group)
	if err != nil {
		return fmt.Errorf("LeaveGroup: error querying database: %w", err)
	}
	if if_group == 0 {
		return fmt.Errorf("LeaveGroup: you don't belong to this group %w", err)
	}

	// rimuovo l'utente dai membri di quella chat
	_, err = db.c.Exec("DELETE FROM membri WHERE id_utenti=$1 AND id_chat=$2", cs, idChat)
	if err != nil {
		return fmt.Errorf("LeaveGroup: error delete user  in table membri: %w", err)
	}

	// vedo se sono presenti altri membri nel gruppo, se cosi non fosse, lo elimino
	err = db.c.QueryRow("SELECT COUNT(id_utenti) FROM membri WHERE id_CHAT=$2", idChat).Scan(&if_group)
	if err != nil {
		return fmt.Errorf("LeaveGroup: error querying database: %w", err)
	}

	// se non sono presenti altri partecipanti allora elimino il gruppo
	if if_group == 0 {
		_, err = db.c.Exec("DELETE FROM chat WHERE id=$2", idChat)
		if err != nil {
			return fmt.Errorf("LeaveGroup: error delete user in table chat: %w", err)
		}
	}

	return nil
}

// imposta il nome del gruppo
func (db *appdbimpl) SetGroupName(cs int, idGroup int, nameGroup string) error {

	var righe int64
	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("SetGroupName: error in authentication: %w", err)
	}

	// vede se l'utente loggato e' presente nel gruppo
	err = db.c.QueryRow("SELECT count(u.id) FROM chat c  JOIN membri m ON m.id_chat=c.id JOIN utenti u ON u.id=m.id_utenti WHERE u.id=$1 AND c.id=$2", cs, idGroup).Scan(&righe)
	if err != nil {
		return fmt.Errorf("SetGroupName: error querying database: %w", err)
	}

	// se non sono state trovate righe allora l'utente loggato non e' nel gruppo
	if righe == 0 {
		return fmt.Errorf("SetGroupName: error you can't change a name of a group that you don't partecipate: %w", err)
	}

	// se l'utente loggato e' nel gruppo allora cambia il nome a quello inserito
	_, err = db.c.Exec("UPDATE chat SET nome=$1 WHERE id=$2", nameGroup, idGroup)
	if err != nil {
		return fmt.Errorf("SetGroupName error: database UPDATE not successful: %w", err)
	}
	return nil

}

// imposta la foto del gruppo
func (db *appdbimpl) SetGroupPhoto(cs int, idGroup int, photoGroup []byte) error {
	var righe int64

	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("SetGroupPhoto: error in authentication: %w", err)
	}

	// vede se l'utente loggato e' presente nel gruppo
	err = db.c.QueryRow("SELECT count(u.id) FROM chat c  JOIN membri m ON m.id_chat=c.id JOIN utenti u ON u.id=m.id_utenti WHERE u.id=$1 AND c.id=$2", cs, idGroup).Scan(&righe)
	if err != nil {
		return fmt.Errorf("SetGroupPhoto: error querying database: %w", err)
	}

	// se non sono state trovate righe allora l'utente loggato non e' nel gruppo
	if righe == 0 {
		return fmt.Errorf("SetGroupPhoto: error you can't change a name of a group that you don't partecipate: %w", err)
	}

	// se l'utente loggato e' nel gruppo allora cambia la foto a quella inserita
	_, err = db.c.Exec("UPDATE chat SET propic=$1 WHERE id=$2", photoGroup, idGroup)
	if err != nil {
		return fmt.Errorf("SetGroupPhoto error: database UPDATE not successful: %w", err)
	}

	return nil
}

// aggiunge uno o piu' utenti ad un gruppo
func (db *appdbimpl) AddToGroup(cs int, idGroup int, membri []int) error {
	var righe int64

	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("AddToGroup: error in authentication: %w", err)
	}

	// vede se l'utente loggato e' presente nel gruppo
	err = db.c.QueryRow("SELECT count(u.id) FROM chat c JOIN membri m ON m.id_chat=c.id JOIN utenti u ON u.id=m.id_utenti WHERE u.id=$1 AND c.id=$2", cs, idGroup).Scan(&righe)
	if err != nil {
		return fmt.Errorf("AddToGroup: error querying database: %w", err)
	}

	// se non sono state trovate righe allora l'utente loggato non e' nel gruppo
	if righe == 0 {
		return fmt.Errorf("AddToGroup: error you can't add user to a group that you don't partecipate: %w", err)
	}

	// altrimenti per ogni membro viene eseguita la query che aggiunge il membro se gia non e' presente nella chat
	for _, c := range membri {

		_, err := db.c.Exec("INSERT INTO membri (id_utenti, id_chat) SELECT $1, $2 WHERE NOT EXISTS (SELECT 1 FROM membri WHERE id_utenti = $1 AND id_chat = $2);", c, idGroup)
		if err != nil {
			return fmt.Errorf("AddToGroup: error insert id_chat and id_utenti in table membri: %w", err)
		}

	}

	return nil
}
