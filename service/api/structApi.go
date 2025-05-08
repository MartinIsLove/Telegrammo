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
type Group struct {
	NomeChat string `json:"nome_chat"`
	Propic   []byte `json:"propic"`
	Membri   []int  `json:"membri"`
	IdChat   int    `json:"id_chat"`
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
	IdMess             int       `json:"id_mess"`
	Nome               string    `json:"nome"` // nome mittente
	IdMitt             int       `json:"id_mitt"`
	Testo              string    `json:"testo"`
	Data               time.Time `json:"data"`
	Visual             bool      `json:"visual"`
	Photo              []byte    `json:"photo"`
	IdChat             int       `json:"id_chat"`
	Gruppo             bool      `json:"gruppo"`
	Emoji              []int     `json:"emoji"` // l'ordine e' tacchi, cuore, ditoSu, ok, manicure
	MyEmoji            string    `json:"myEmoji"`
	ForwardUsername    string    `json:"forward_username"` // username di colui che ha inviato il messaggio da forwardare
	ForwardId          int       `json:"forward_id"`       // id di colui che ha inviato il messaggio da forwardare
	ForwardDate        time.Time `json:"forward_date"`     // data forward
	ForwardIdMit       int       `json:"forward_id_mit"`   // id di colui che ha fatto il forward
	ForwardUsernameMit string    `json:"forward_user_mit"` // username di colui che ha fatto il forward
	IdForward          int       `json:"id_forward"`       // id del forward
	Idreply            int       `json:"id_reply"`         // l'id del messaggio di chat a cui fa riferimento il reply
	TestoReply         string    `json:"testo_reply"`      // il testo del messaggio del reply se presente
	IdMitReply         int       `json:"id_mit_reply"`     // l'id del mittente del reply
	MitReply           string    `json:"mit_reply"`        // l'username del mittente del reply
	PhotoReply         []byte    `json:"photo_reply"`      // la photo a cui potrebbe far riferimento il reply
}

type NomeChat struct {
	Gruppo  bool   `json:"gruppo"`
	Nome    string `json:"nome"`
	Message []Mess `json:"message"`
}
type Comment struct {
	Id_mes  int    `json:"id_mes"`
	Emoji   string `json:"emoji"`
	Id_chat int    `json:"id_chat"`
}

type RequestData struct {
	IdChat    []int `json:"id_chat"`
	IdMes     int   `json:"id_mes"`
	IdForward int   `json:"id_for"`
}

func NewMess(messaggio database.MessDb) Mess {
	return Mess{
		IdMess:             messaggio.IdMess,
		Nome:               messaggio.Nome,
		IdMitt:             messaggio.IdMitt,
		Testo:              messaggio.Testo,
		Data:               messaggio.Data,
		Visual:             messaggio.Visual,
		Photo:              messaggio.Photo,
		IdChat:             messaggio.IdChat,
		Gruppo:             messaggio.Gruppo,
		Emoji:              messaggio.Emoji,
		MyEmoji:            messaggio.MyEmoji,
		ForwardUsername:    messaggio.ForwardUsername,
		ForwardId:          messaggio.ForwardId,
		ForwardDate:        messaggio.ForwardDate,
		ForwardIdMit:       messaggio.ForwardIdMit,
		ForwardUsernameMit: messaggio.ForwardUsernameMit,
		IdForward:          messaggio.IdForward,
		Idreply:            messaggio.Idreply,
		IdMitReply:         messaggio.IdMitReply,
		TestoReply:         messaggio.TestoReply,
		MitReply:           messaggio.MitReply,
		PhotoReply:         messaggio.PhotoReply,
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
