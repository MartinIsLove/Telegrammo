package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var richiesta Mess
	cs, err := rt.AuthenticationApi(r)
	if err != nil {
		http.Error(w, "error: authentication user sendMessage"+err.Error(), http.StatusUnauthorized)
		return
	}
	// id_chat_tmp := ps.ByName("idChat")
	// id_chat, err := strconv.Atoi(id_chat_tmp)
	// if err != nil {
	// 	http.Error(w, "error: conversion id_chat_tmp (string) to id_chat (int)", http.StatusInternalServerError)
	// 	return
	// }

	if err := json.NewDecoder(r.Body).Decode(&richiesta); err != nil {
		http.Error(w, "sendMessage:"+err.Error(), http.StatusBadRequest)
		return
	}
	if len(richiesta.Testo) < 1 {
		http.Error(w, "messaggio troppo corto", http.StatusBadRequest)
		return
	}
	erro := rt.db.SendMessage(cs, richiesta.IdChat, richiesta.Testo, richiesta.Photo)
	if erro != nil && strings.HasPrefix(erro.Error(), "error in authentication SendMessage:") {
		http.Error(w, erro.Error(), http.StatusUnauthorized)
		return
	}
	if erro != nil {
		http.Error(w, erro.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
