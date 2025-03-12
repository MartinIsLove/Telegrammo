package api

import "github.com/MartinIsLove/Telegrammo/service/database"

type Utente struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
	Propic   []byte `json:"propic"`
}

type Chat struct {
	Id     int    `json:"id"`
	Nome   string `json:"nome"`
	Propic []byte `json:"propic"`
	Gruppo bool   `json:"gruppo"`
}

func NewUser(user database.UtenteDb) Utente {
	return Utente{
		Id:       user.Id,
		Username: user.Username,
		Propic:   user.Propic,
	}
}
