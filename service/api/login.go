package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var richiestaLogin Utente

	err := json.NewDecoder(r.Body).Decode(&richiestaLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// fmt.Println(richiestaLogin)
	var rispostaLogin Utente

	rispostaLogin.Id, err = rt.db.DoLogin(richiestaLogin.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
