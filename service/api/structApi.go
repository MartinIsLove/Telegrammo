package api

import (
	"time"

	"github.com/MartinIsLove/Telegrammo/service/database"
)

type Utente struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
	Propic   []byte `json:"propic"`
}

type ChatUtente struct {
	IdChat   int       `json:"id_chat"`
	Nome     string    `json:"nome"`
	Gruppo   bool      `json:"gruppo"`
	Username string    `json:"username"`
	Id       int       `json:"id_utenti"`
	Propic   []byte    `json:"propic"`
	LastMSg  string    `json:"lastmsg"`
	Data     time.Time `json:"data"`
}

type Mess struct {
	IdMess int       `json:"id_mess"`
	Nome   string    `json:"nome"`
	IdMitt int       `json:"id_mitt"`
	Testo  string    `json:"testo"`
	Data   time.Time `json:"data"`
	Visual bool      `json:"visual"`
	Photo  []byte    `json:"photo"`
}

func NewMess(Messaggio database.MessDb) Mess {
	return Mess{
		IdMess: Messaggio.IdMess,
		Nome:   Messaggio.Nome,
		IdMitt: Messaggio.IdMitt,
		Testo:  Messaggio.Testo,
		Data:   Messaggio.Data,
		Visual: Messaggio.Visual,
	}
}
func NewChatUtente(chatUtente database.ChatUtenteDb) ChatUtente {
	return ChatUtente{
		IdChat:   chatUtente.IdChat,
		Nome:     chatUtente.Nome,
		Gruppo:   chatUtente.Gruppo,
		Username: chatUtente.Username,
		Id:       chatUtente.Id,
		Propic:   chatUtente.Propic,
		LastMSg:  chatUtente.LastMSg,
		Data:     chatUtente.Data,
	}
}
func NewUser(user database.UtenteDb) Utente {
	return Utente{
		Id:       user.Id,
		Username: user.Username,
		Propic:   user.Propic,
	}
}
