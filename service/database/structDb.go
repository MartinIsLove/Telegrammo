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
type MessDb struct {
	IdMess             int       `json:"id_mess"`
	Nome               string    `json:"nome"`
	IdMitt             int       `json:"id_mitt"`
	Testo              string    `json:"testo"`
	Data               time.Time `json:"data"`
	Visual             bool      `json:"visual"`
	Photo              []byte    `json:"photo"`
	IdChat             int       `json:"id_chat"`
	Gruppo             bool      `json:"grupo"`
	Emoji              []int     `json:"emoji"` // l'ordine e' tacchi, cuore, ditoSu, ok, manicure
	MyEmoji            string    `json:"myemoji"`
	ForwardUsername    string    `json:"forward_username"` // username di colui che ha inviato il messaggio da forwardare
	ForwardId          int       `json:"forward_id"`       // id di colui che ha inviato il messaggio da forwardare
	ForwardDate        time.Time `json:"forward_date"`     // data forward
	ForwardIdMit       int       `json:"forward_id_mit"`   // id di colui che ha fatto il forward
	ForwardUsernameMit string    `json:"forward_user_mit"` // username di colui che ha fatto il forward
	IdForward          int       `json:"id_forward"`       //id del forward
}
