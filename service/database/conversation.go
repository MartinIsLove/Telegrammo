package database

import (
	"fmt"
)

func (db *appdbimpl) IsChatDuplicated(cs int, id int) (bool, error) {
	var num_righe int64

	err := db.c.QueryRow("SELECT COUNT(chat.id) AS righe FROM chat JOIN membri ON chat.id=membri.id_chat WHERE chat.gruppo=FALSE AND membri.id_utenti IN ($1 , $2) GROUP BY chat.id HAVING COUNT (DISTINCT membri.id_utenti)=2", cs, id).Scan(&num_righe)
	if err != nil {
		return true, fmt.Errorf("user: error checking chat duplicated in database: %w", err)
	}
	if num_righe == 0 {
		return false, nil
	}
	return true, nil
}
func (db *appdbimpl) CreateChat(cs int, id int) error {

}
