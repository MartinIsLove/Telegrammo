package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
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
func (rt *_router) checknames(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var richiesta []Utente // cerca tutti gli utenti che iniziano per una data stringa
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
	if len(str) < 1 && len(str) > 16 {
		http.Error(w, "too short message", http.StatusBadRequest)
	}

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
		http.Error(w, "error: authentication user getMyConversations"+err.Error(), http.StatusUnauthorized)
		return
	}
	tmp, err := rt.db.GetMyConversations(cs)
	if err != nil && strings.HasPrefix(err.Error(), "error in authentication GetConversations:") {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && strings.HasPrefix(err.Error(), "GetConversations: no chats or groups find:") {

	} else if err != nil {
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
func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	cs, err := rt.AuthenticationApi(r)
	var messaggi []Mess

	if err != nil {
		http.Error(w, "error: authentication user getConversation"+err.Error(), http.StatusUnauthorized)
		return
	}
	// var id_chat_tmp string
	id_chat_tmp := ps.ByName("idChat")
	id_chat, err := strconv.Atoi(id_chat_tmp)
	if err != nil {
		http.Error(w, "error getConversation: conversion id_chat_tmp (string) to id_chat (int)", http.StatusBadRequest)
	}

	gp, str, tmp, err := rt.db.GetConversation(cs, id_chat)
	if err != nil && strings.HasPrefix(err.Error(), "error in authentication GetConversation:") {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && strings.HasPrefix(err.Error(), "GetConversation: no messages found:") {

	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	messaggi = make([]Mess, len(tmp))
	for i, user := range tmp {
		messaggi[i] = NewMess(user)
	}

	var Nome NomeChat
	Nome.Gruppo = gp
	Nome.Nome = str
	Nome.Message = messaggi
	w.Header().Set("content-type", "application/json")

	json, err := json.Marshal(Nome)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(json)

	w.WriteHeader(http.StatusOK)

}
func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var richiesta Group

	cs, err := rt.AuthenticationApi(r)

	if err != nil {
		http.Error(w, "error: authentication user createChat"+err.Error(), http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&richiesta); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = rt.db.CreateGroup(cs, richiesta.NomeChat, richiesta.Propic, richiesta.Membri)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cs, err := rt.AuthenticationApi(r)

	if err != nil {
		http.Error(w, "error: authentication user createChat"+err.Error(), http.StatusUnauthorized)
		return
	}

	id_chat_tmp := ps.ByName("idChat")
	id_chat, err := strconv.Atoi(id_chat_tmp)
	if err != nil {
		http.Error(w, "error leaveGroup: conversion id_chat_tmp (string) to id_chat (int)", http.StatusBadRequest)
		return
	}

	err = rt.db.LeaveGroup(cs, id_chat)
	if err != nil {
		http.Error(w, "error leaveGroup: conversion id_chat_tmp (string) to id_chat (int)", http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusOK)
}
