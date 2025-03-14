package database

import "time"

type UtenteDb struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
	Propic   []byte `json:"propic"`
}
type ChatUtenteDb struct {
	IdChat   int       `json:"id_chat"`
	Nome     string    `json:"nome"`
	Gruppo   bool      `json:"gruppo"`
	Username string    `json:"username"`
	Id       int       `json:"id_utenti"`
	Propic   []byte    `json:"propic"`
	LastMSg  string    `json:"lastmsg"`
	Data     time.Time `json:"data"`
}
