package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) createChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var richiesta Utente

	cs, err := rt.AuthenticationApi(r)

	if err != nil {
		http.Error(w, "error: authentication user createChat"+err.Error(), http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&richiesta); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = rt.db.CreateChat(cs, richiesta.Id)

	if err != nil && strings.HasPrefix(err.Error(), "chat: error this chat already exist:") {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil && strings.HasPrefix(err.Error(), "error in authentication CreateChat:") {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (rt *_router) checknames(w http.ResponseWriter, r *http.Request, ps httprouter.Params) { // cerca tutti gli utenti che iniziano per una data stringa
	var richiesta []Utente
	var str string
	cs, err := rt.AuthenticationApi(r)

	if err != nil {
		http.Error(w, "error: authentication user checkname"+err.Error(), http.StatusUnauthorized)
		return
	}
	str = ps.ByName("username")
	if len(str) == 0 {
		http.Error(w, "error: username too short", http.StatusBadRequest)
		return
	}
	// if err := json.NewDecoder(r.Body).Decode(&str); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	tmp, err := rt.db.CheckNames(cs, str)

	if err != nil && strings.HasPrefix(err.Error(), "error in authentication Checknames:") {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && strings.HasPrefix(err.Error(), "chat: no user found:") {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error in database"+err.Error(), http.StatusInternalServerError)
		return
	}

	richiesta = make([]Utente, len(tmp))
	for i, user := range tmp {
		richiesta[i] = NewUser(user)
	}

	w.Header().Set("content-type", "application/json")
	json, err := json.Marshal(richiesta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(json)

	w.WriteHeader(http.StatusOK)
}

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var richiesta []ChatUtente
	cs, err := rt.AuthenticationApi(r)
	if err != nil {
		http.Error(w, "error: authentication user checkname"+err.Error(), http.StatusUnauthorized)
		return
	}
	tmp, err := rt.db.GetMyConversations(cs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	richiesta = make([]ChatUtente, len(tmp))
	for i, user := range tmp {
		richiesta[i] = NewChatUtente(user)
	}

	sort.Slice(richiesta, func(i, j int) bool {
		if richiesta[i].Data.IsZero() {
			return false
		}
		if richiesta[j].Data.IsZero() {
			return true
		}
		return richiesta[i].Data.After(richiesta[j].Data)
	})

	w.Header().Set("content-type", "application/json")
	json, err := json.Marshal(richiesta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(json)

	w.WriteHeader(http.StatusOK)
	// fmt.Println(err.Error())
}

// func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

// 	cs, err := rt.AuthenticationApi(r)
// 	var messaggi []Mess

// 	if err != nil {
// 		http.Error(w, "error: authentication user checkname"+err.Error(), http.StatusUnauthorized)
// 		return
// 	}
// 	// var id_chat_tmp string
// 	id_chat_tmp := ps.ByName("idChat")
// 	id_chat, err := strconv.Atoi(id_chat_tmp)

// 	tmp, err := rt.db.GetConversation(cs, id_chat)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	messaggi = make([]Mess, len(tmp))
// 	for i, user := range tmp {
// 		messaggi[i] = NewMess(user)
// 	}
// }
