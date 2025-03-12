package database

type UtenteDb struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
	Propic   []byte `json:"propic"`
}
