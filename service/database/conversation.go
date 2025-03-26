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
func (db *appdbimpl) CreateChat(cs int, id int) (int, error) {

	_, err := db.Authentication(cs)
	if err != nil {
		return -1, fmt.Errorf("error in authentication CreateChat: %w", err)
	}

	Isduplicated, err := db.IsChatDuplicated(cs, id)

	if err != nil {
		return -1, err
	}

	if !Isduplicated {
		str_chat, err := db.c.Exec("INSERT INTO chat (gruppo) VALUES (FALSE)")
		if err != nil {
			return -1, fmt.Errorf("chat: error insert chat in table chat: %w", err)
		}

		id_chat, err := str_chat.LastInsertId()
		if err != nil {
			return -1, fmt.Errorf("chat: error catch number of rows from query: %w", err)
		}

		_, err = db.c.Exec("INSERT INTO membri (id_chat, id_utenti) VALUES  ($1, $2), ($1, $3)", id_chat, cs, id)
		if err != nil {
			return -1, fmt.Errorf("chat: error insert chat users in membri : %w", err)
		}
		return int(id_chat), nil
	} else {
		return -1, err
	}

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
		return utenti, fmt.Errorf("checkNames: error querying users: %w", err)
	}

	var cont int

	defer rows.Close()

	for rows.Next() {
		cont++
		var utente UtenteDb
		if err := rows.Scan(&utente.Id, &utente.Username, &utente.Propic); err != nil {
			return utenti, fmt.Errorf("checkNames: error scanning user: %w", err)
		}
		utenti = append(utenti, utente)
	}
	if cont == 0 {
		return utenti, fmt.Errorf("checkNames: no user found: %w", err)
	}
	if err := rows.Err(); err != nil {
		return utenti, fmt.Errorf("checkNames: error iterating over users: %w", err)
	}

	// ----------------------------------------------

	return utenti, nil
}
func (db *appdbimpl) CheckChatNames(cs int, toFind string) ([]ChatUtenteDb, error) {

	var chat []ChatUtenteDb
	_, err := db.Authentication(cs)
	if err != nil {
		return []ChatUtenteDb{}, fmt.Errorf("error in authentication GetConversations: %w", err)
	}
	fmt.Println(cs, toFind)

	// la query sotto ritorna gli id degli utenti con cui ha la chat l'utente connesso, che non siano gruppi e l'id della chat
	rows, err1 := db.c.Query("SELECT u.username, u.propic, m.id_utenti, c.id FROM chat c JOIN membri m ON c.id=m.id_chat JOIN utenti u ON u.id = m.id_utenti WHERE (u.username LIKE $1 || '%') AND m.id_utenti!=$2 AND c.id IN(SELECT  c.id from chat c JOIN membri m ON c.id=m.id_chat WHERE m.id_utenti=$2 AND gruppo=0) AND gruppo=0;", toFind, cs)
	if err1 != nil && !errors.Is(err1, sql.ErrNoRows) {
		return []ChatUtenteDb{}, fmt.Errorf("GetConversations: error querying users: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var c ChatUtenteDb
		if err := rows.Scan(&c.Nome, &c.Propic, &c.Id, &c.IdChat); err != nil {
			return []ChatUtenteDb{}, fmt.Errorf(" Getconversations: error scanning user: %w", err)
		}
		fmt.Println(c, "dopca")
		chat = append(chat, c)
	}
	fmt.Println(chat)
	if err := rows.Err(); err != nil {
		return []ChatUtenteDb{}, fmt.Errorf("GetConversations: error iterating over users: %w", err)
	}

	// questa query ritorna tutti i dati della join tra membri e chat dove l'utente appartiene al gruppo
	rows2, err2 := db.c.Query("SELECT  c.id AS id, c.propic, c.nome, c.gruppo from chat c JOIN membri m ON c.id=m.id_chat JOIN utenti u ON u.id=m.id_utenti WHERE (c.nome LIKE $1 || '%') AND m.id_utenti=$2 AND gruppo=1;", toFind, cs)
	if err2 != nil && !errors.Is(err2, sql.ErrNoRows) {
		return []ChatUtenteDb{}, fmt.Errorf("GetConversations: error querying users: %w", err)
	}

	defer rows2.Close()

	for rows2.Next() {
		var c ChatUtenteDb
		if err := rows2.Scan(&c.IdChat, &c.Propic, &c.Nome, &c.Gruppo); err != nil {
			return []ChatUtenteDb{}, fmt.Errorf("GetConversations: error scanning user: %w", err)
		}
		chat = append(chat, c)
	}
	if err := rows2.Err(); err != nil {
		return []ChatUtenteDb{}, fmt.Errorf("GetConversations: error iterating over users: %w", err)
	}

	if errors.Is(err1, sql.ErrNoRows) && errors.Is(err2, sql.ErrNoRows) {
		return []ChatUtenteDb{}, fmt.Errorf("GetConversations: no chats or groups find: %w", err)

	}

	return chat, nil
}
func (db *appdbimpl) GetMyConversations(cs int) ([]ChatUtenteDb, error) {
	var chat []ChatUtenteDb
	_, err := db.Authentication(cs)
	if err != nil {
		return []ChatUtenteDb{}, fmt.Errorf("error in authentication GetConversations: %w", err)
	}
	// la query sotto ritorna gli id degli utenti con cui ha la chat l'utente connesso, che non siano gruppi e l'id della chat
	rows, err1 := db.c.Query("SELECT u.username, u.propic, m.id_utenti, c.id FROM chat c JOIN membri m ON c.id=m.id_chat JOIN utenti u ON u.id = m.id_utenti WHERE m.id_utenti!=$1 AND c.id IN(SELECT  c.id from chat c JOIN membri m ON c.id=m.id_chat WHERE m.id_utenti=$1 AND gruppo=0) AND gruppo=0;", cs)
	if err1 != nil && !errors.Is(err1, sql.ErrNoRows) {
		return []ChatUtenteDb{}, fmt.Errorf("GetConversations: error querying users: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var c ChatUtenteDb
		if err := rows.Scan(&c.Nome, &c.Propic, &c.Id, &c.IdChat); err != nil {
			return []ChatUtenteDb{}, fmt.Errorf(" Getconversations: error scanning user: %w", err)
		}
		chat = append(chat, c)
	}
	if err := rows.Err(); err != nil {
		return []ChatUtenteDb{}, fmt.Errorf("GetConversations: error iterating over users: %w", err)
	}

	// questa query ritorna tutti i dati della join tra membri e chat dove l'utente appartiene al gruppo
	rows2, err2 := db.c.Query("SELECT  c.id AS id, c.propic, c.nome, c.gruppo from chat c JOIN membri m ON c.id=m.id_chat JOIN utenti u ON u.id=m.id_utenti WHERE m.id_utenti=$1 AND gruppo=1;", cs)
	if err2 != nil && !errors.Is(err2, sql.ErrNoRows) {
		return []ChatUtenteDb{}, fmt.Errorf("GetConversations: error querying users: %w", err)
	}

	defer rows2.Close()

	for rows2.Next() {
		var c ChatUtenteDb
		if err := rows2.Scan(&c.IdChat, &c.Propic, &c.Nome, &c.Gruppo); err != nil {
			return []ChatUtenteDb{}, fmt.Errorf("GetConversations: error scanning user: %w", err)
		}
		chat = append(chat, c)
	}
	if err := rows2.Err(); err != nil {
		return []ChatUtenteDb{}, fmt.Errorf("GetConversations: error iterating over users: %w", err)
	}

	if errors.Is(err1, sql.ErrNoRows) && errors.Is(err2, sql.ErrNoRows) {
		return []ChatUtenteDb{}, fmt.Errorf("GetConversations: no chats or groups find: %w", err)

	}

	for i := range chat {
		var lastMsg, username string
		var id int
		var tmp time.Time
		c := &chat[i]
		rows3 := db.c.QueryRow("SELECT m.testo, m.mittente, m.data FROM chat c JOIN messaggi_di_chat mdc ON c.id=mdc.id_chat JOIN messaggi m ON mdc.id_messaggio=m.id JOIN membri me ON me.id_chat=c.id JOIN utenti u ON u.id=me.id_utenti WHERE c.id=$1 ORDER BY m.data DESC LIMIT 1;", c.IdChat)
		err := rows3.Scan(&lastMsg, &id, &tmp)
		if err == sql.ErrNoRows {

			c.LastMSg = ""
			c.Data = time.Time{}

		} else if err != nil {
			return []ChatUtenteDb{}, fmt.Errorf("GetConversations: error querying users: %w", err)

		} else {
			if len(lastMsg) > 100 {
				c.LastMSg = lastMsg[:100] + "..."
			} else {
				c.LastMSg = lastMsg
			}
			c.Id = id
			c.Data = tmp
		}
		rows4 := db.c.QueryRow("SELECT username FROM utenti WHERE id=$1 LIMIT 1", id)
		err = rows4.Scan(&username)
		if err == sql.ErrNoRows {

		} else if err != nil {
			return []ChatUtenteDb{}, fmt.Errorf("GetConversations: error querying users: %w", err)

		}
		c.Username = username
	}
	// for _, c := range chat {
	// 	fmt.Printf("il nome della chat e': %s Username di chi ha inviato l'ultimo messaggio: %s, Propic: %s, UserId di chi ha inviato l'ultimo messaggio: %d, ChatId: %d , lastmsg: %s, data: %s, e' un gruppo: %t \n", c.Nome, c.Username, c.Propic, c.Id, c.IdChat, c.LastMSg, c.Data.Format(time.RFC3339), c.Gruppo)
	// }

	return chat, nil
}
func (db *appdbimpl) GetConversation(cs int, id_chat int) (bool, string, []MessDb, error) {
	var mess []MessDb
	_, err := db.Authentication(cs)
	if err != nil {
		return false, "", []MessDb{}, fmt.Errorf("error in authentication GetConversation: %w", err)
	}

	var group bool

	err = db.c.QueryRow("SELECT gruppo FROM chat WHERE id=$1 LIMIT 1;", id_chat).Scan(&group)

	if err != nil {
		return false, "", []MessDb{}, fmt.Errorf("GetConversation: error querying users: %w", err)
	}
	var nomeChat string
	if group {
		err = db.c.QueryRow("SELECT nome FROM chat WHERE id=$1 LIMIT 1;", id_chat).Scan(&nomeChat)
		if err != nil {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: error querying users: %w", err)
		}
	} else {
		err = db.c.QueryRow("SELECT utenti.username FROM chat c JOIN membri ON c.id = membri.id_chat JOIN utenti ON membri.id_utenti = utenti.id WHERE c.id=$1 AND utenti.id != $2 LIMIT 1;", id_chat, cs).Scan(&nomeChat)
		if err != nil {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: error querying users: %w", err)
		}
	}
	// prova a fare una unica query qua sopra che prende tutto qui
	rows, err := db.c.Query("SELECT c.gruppo, m.testo, m.mittente, u.username, m.data, m.image, m.id,d.id FROM messaggi m JOIN messaggi_di_chat d ON d.id_messaggio=m.id JOIN chat c ON c.id=d.id_chat JOIN utenti u ON u.id=m.mittente WHERE c.id=$1 ORDER BY m.data;", id_chat)
	if errors.Is(err, sql.ErrNoRows) {
		return false, "", []MessDb{}, fmt.Errorf("GetConversation: no messages found: %w", err)
	}
	if err != nil {
		return false, "", []MessDb{}, fmt.Errorf("GetConversation: error querying users: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var c MessDb
		var i int
		if err := rows.Scan(&c.Gruppo, &c.Testo, &c.IdMitt, &c.Nome, &c.Data, &c.Photo, &c.IdMess, &i); err != nil {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: error scanning user: %w", err)
		}

		var forwarderUsername string
		var forwarderId int
		var forwardDate time.Time
		var forwardIdMit int
		var forwarderUsernameMit string
		var cont int = 0

		row2, err := db.c.Query("SELECT u.username, u.id, m.forward_date, m.id_forw_mit, t.username FROM utenti u JOIN messaggi_di_chat m ON m.id_forward = u.id JOIN utenti t ON m.id_forw_mit=t.id WHERE m.id_messaggio = $1 AND m.id_forward NOT NULL AND m.id=$2", c.IdMess, i) //.Scan(&forwarderUsername, &forwarderId, &forwardDate, &forwardIdMit, &forwarderUsernameMit)
		if errors.Is(err, sql.ErrNoRows) {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: no messages found: %w", err)
		}
		if err != nil {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: error querying users: %w", err)
		}

		defer row2.Close()

		for row2.Next() {
			if err := row2.Scan(&forwarderUsername, &forwarderId, &forwardDate, &forwardIdMit, &forwarderUsernameMit); err != nil {
				return false, "", []MessDb{}, fmt.Errorf("GetConversation: error scanning user: %w", err)
			}

			c.ForwardUsername = forwarderUsername
			c.ForwardId = forwarderId
			c.ForwardDate = forwardDate
			c.ForwardIdMit = forwardIdMit
			c.ForwardUsernameMit = forwarderUsernameMit
			mess = append(mess, c)
			cont++
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
		return false, "", []MessDb{}, fmt.Errorf("GetConversation: error iterating over users: %w", err)
	}

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

	for _, m := range mess {
		var emoji string

		err = db.c.QueryRow("SELECT emoji FROM emoticon WHERE id_utente=$1 AND id_messaggio=$2", cs, m.IdMess).Scan(&emoji)
		if err == sql.ErrNoRows {
			emoji = ""
		} else if err != nil {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: error querying users: %w", err)
		}

		var counts [5]int
		err = db.c.QueryRow("SELECT COALESCE(SUM(CASE WHEN emoji = '👠' AND id_messaggio=$1 THEN 1 ELSE 0 END),0) AS tacchi, COALESCE(SUM(CASE WHEN emoji = '❤️' AND id_messaggio=$1 THEN 1 ELSE 0 END),0) AS cuore, COALESCE(SUM(CASE WHEN emoji = '👍🏻' AND id_messaggio=$1 THEN 1 ELSE 0 END),0) AS dito_in_su, COALESCE(SUM(CASE WHEN emoji = '👌🏻' AND id_messaggio=$1 THEN 1 ELSE 0 END),0) AS ok, COALESCE(SUM(CASE WHEN emoji = '💅' AND id_messaggio=$1 THEN 1 ELSE 0 END),0) AS manicure FROM emoticon;", m.IdMess).Scan(&counts[0], &counts[1], &counts[2], &counts[3], &counts[4])
		if err != nil {
			return false, "", []MessDb{}, fmt.Errorf("GetConversation: error querying users: %w", err)
		}

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
func (db *appdbimpl) CreateGroup(cs int, nome string, propic []byte, membri []int) (int, error) {
	_, err := db.Authentication(cs)
	if err != nil {
		return -1, fmt.Errorf("error in authentication CreateGroup: %w", err)
	}
	noPhotoPath := filepath.Join("/workspace/webui/src/assets/", "NoPhoto.png")
	noPhotoBytes, err := os.ReadFile(noPhotoPath)
	if err != nil {
		return -1, fmt.Errorf("createGroup: error reading noPhoto.png: %w", err)
	}
	if propic == nil {
		propic = noPhotoBytes
	}

	str_group, err := db.c.Exec("INSERT INTO chat (nome, propic, gruppo) VALUES ($1, $2, TRUE)", nome, propic)
	if err != nil {
		return -1, fmt.Errorf("group: error insert group in table chat: %w", err)
	}

	id_group, err := str_group.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("chat: error catch number of rows from query: %w", err)
	}
	_, err = db.c.Exec("INSERT INTO membri (id_utenti, id_chat) VALUES ($1, $2)", cs, id_group)
	if err != nil {
		return -1, fmt.Errorf("group: error insert group in table chat: %w", err)
	}
	for _, m := range membri {
		value := db.UserExist(m)
		if value {
			_, err := db.c.Exec("INSERT INTO membri (id_utenti, id_chat) VALUES ($1, $2)", m, id_group)
			if err != nil {
				return -1, fmt.Errorf("group: error insert group in table chat: %w", err)
			}
		}
	}

	return int(id_group), nil
}
func (db *appdbimpl) LeaveGroup(cs int, idChat int) error {

	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("error in authentication CreateGroup: %w", err)
	}

	_, err = db.c.Exec("DELETE FROM membri WHERE id_utenti=$1 AND id_chat=$2", cs, idChat)
	if err != nil {
		return fmt.Errorf("group: error leave user by a group in table membri: %w", err)
	}

	return nil
}
func (db *appdbimpl) SetGroupName(cs int, idGroup int, nameGroup string) error {

	var righe int64
	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("SetGroupName: error in authentication: %w", err)
	}

	err = db.c.QueryRow("SELECT count(u.id) FROM chat c  JOIN membri m ON m.id_chat=c.id JOIN utenti u ON u.id=m.id_utenti WHERE u.id=$1 AND c.id=$2", cs, idGroup).Scan(&righe)
	if err != nil {
		return fmt.Errorf("SetGroupName: error querying database: %w", err)
	}
	if righe == 0 {
		return fmt.Errorf("you can't change a name of a group that you don't partecipate: %w", err)
	}

	_, err = db.c.Exec("UPDATE chat SET nome=$1 WHERE id=$2", nameGroup, idGroup)
	if err != nil {
		return fmt.Errorf("SetGroupName error: database UPDATE not successful: %w", err)
	}
	return nil

}
func (db *appdbimpl) SetGroupPhoto(cs int, idGroup int, photoGroup []byte) error {
	var righe int64

	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("SetGroupPhoto: error in authentication: %w", err)
	}

	err = db.c.QueryRow("SELECT count(u.id) FROM chat c  JOIN membri m ON m.id_chat=c.id JOIN utenti u ON u.id=m.id_utenti WHERE u.id=$1 AND c.id=$2", cs, idGroup).Scan(&righe)
	if err != nil {
		return fmt.Errorf("SetGroupPhoto: error querying database: %w", err)
	}
	if righe == 0 {
		return fmt.Errorf("you can't change a name of a group that you don't partecipate: %w", err)
	}

	_, err = db.c.Exec("UPDATE chat SET propic=$1 WHERE id=$2", photoGroup, idGroup)
	if err != nil {
		return fmt.Errorf("SetGroupPhoto error: database UPDATE not successful: %w", err)
	}

	return nil
}
func (db *appdbimpl) AddToGroup(cs int, idGroup int, membri []int) error {
	var righe int64

	_, err := db.Authentication(cs)
	if err != nil {
		return fmt.Errorf("AddToGroup: error in authentication: %w", err)
	}

	err = db.c.QueryRow("SELECT count(u.id) FROM chat c JOIN membri m ON m.id_chat=c.id JOIN utenti u ON u.id=m.id_utenti WHERE u.id=$1 AND c.id=$2", cs, idGroup).Scan(&righe)
	if err != nil {
		return fmt.Errorf("AddToGroup: error querying database: %w", err)
	}
	if righe == 0 {
		return fmt.Errorf("you can't add user to a group that you don't partecipate: %w", err)
	}

	for _, c := range membri {

		_, err := db.c.Exec("INSERT INTO membri (id_utenti, id_chat) SELECT $1, $2 WHERE NOT EXISTS (SELECT 1 FROM membri WHERE id_utenti = $1 AND id_chat = $2);", c, idGroup)
		if err != nil {
			return fmt.Errorf("AddToGroup: error insert id_chat and id_utenti in table membri: %w", err)
		}

	}

	return nil
}
