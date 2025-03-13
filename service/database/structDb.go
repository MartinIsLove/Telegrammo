package database

type UtenteDb struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
	Propic   []byte `json:"propic"`
}
type ChatDb struct {
	Id     int    `json:"id"`
	Nome   string `json:"nome"`
	Propic []byte `json:"propic"`
	Gruppo bool   `json:"gruppo"`
}
