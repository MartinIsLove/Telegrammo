package api

import (
	"encoding/json"
	"net/http"
	"unicode"

	"github.com/julienschmidt/httprouter"
)

// verifica se il primo carattere di una stringa e' alfanumerico
func isFirstCharAlphanumeric(s string) bool {
	if len(s) == 0 {
		return false
	}
	r := rune(s[0])
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

// si occupa del login
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var richiestaLogin Utente

	err := json.NewDecoder(r.Body).Decode(&richiestaLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var rispostaLogin Utente

	if len(richiestaLogin.Username) < 1 || len(richiestaLogin.Username) > 16 {
		http.Error(w, "username too short or too long", http.StatusBadRequest)
		return
	}
	if !isFirstCharAlphanumeric(richiestaLogin.Username) {
		http.Error(w, "the first digit of the username must be an alfanumeric character", http.StatusBadRequest)
		return
	}

	rispostaLogin.Id, err = rt.db.DoLogin(richiestaLogin.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	idJson, err := json.Marshal(rispostaLogin.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(idJson)

}
