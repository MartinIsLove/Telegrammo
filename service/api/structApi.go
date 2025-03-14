package api

import "github.com/MartinIsLove/Telegrammo/service/database"

type Utente struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
	Propic   []byte `json:"propic"`
}

type ChatUtente struct {
	IdChat   int    `json:"id_chat"`
	Nome     string `json:"nome"`
	Gruppo   bool   `json:"gruppo"`
	Username string `json:"username"`
	Id       int    `json:"id_utenti"`
	Propic   []byte `json:"propic"`
}

func NewChatUtente(chatUtente database.ChatUtenteDb) ChatUtente {
	return ChatUtente{
		IdChat:   chatUtente.IdChat,
		Nome:     chatUtente.Nome,
		Gruppo:   chatUtente.Gruppo,
		Username: chatUtente.Username,
		Id:       chatUtente.Id,
		Propic:   chatUtente.Propic,
	}
}
func NewUser(user database.UtenteDb) Utente {
	return Utente{
		Id:       user.Id,
		Username: user.Username,
		Propic:   user.Propic,
	}
}
