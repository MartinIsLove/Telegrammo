package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) createChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var richiesta Utente

	cs, err := rt.AuthenticationApi(r)

	if err != nil {
		http.Error(w, "error: authentication user createChat"+err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&richiesta); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = rt.db.CreateChat(cs, richiesta.Id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (rt *_router) checknames(w http.ResponseWriter, r *http.Request, ps httprouter.Params) { // cerca tutti gli utenti che iniziano per una data stringa
	var richiesta []Utente
	var str Utente
	cs, err := rt.AuthenticationApi(r)

	if err != nil {
		http.Error(w, "error: authentication user checkname"+err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&str); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmp, err := rt.db.CheckNames(cs, str.Username)
	if err != nil {
		http.Error(w, "error in database"+err.Error(), http.StatusBadRequest)
		return
	}

	richiesta = make([]Utente, len(tmp))
	for i, user := range tmp {
		richiesta[i] = NewUser(user)
	}
	fmt.Println(richiesta)

}
