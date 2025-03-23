package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

	id_chat, err := rt.db.CreateChat(cs, richiesta.Id)

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

	response := map[string]interface{}{
		"id_chat": id_chat,
	}
	w.Header().Set("content-type", "application/json")
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(json)

	w.WriteHeader(http.StatusOK)
}
func (rt *_router) checknames(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var richiesta []Utente // cerca tutti gli utenti che iniziano per una data stringa
	var str string
	cs, err := rt.AuthenticationApi(r)

	if err != nil {
		http.Error(w, "error: authentication user checkname "+err.Error(), http.StatusUnauthorized)
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
	if len(str) > 16 {
		http.Error(w, "too long message", http.StatusBadRequest)
		return
	}
	if str == "$" {
		w.WriteHeader(http.StatusOK)
		return
	}
	tmp, err := rt.db.CheckNames(cs, str)

	if err != nil && strings.HasPrefix(err.Error(), "error in authentication Checknames:") {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && strings.HasPrefix(err.Error(), "checkNames: no user found:") {
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

	var nome NomeChat
	nome.Gruppo = gp
	nome.Nome = str
	nome.Message = messaggi
	w.Header().Set("content-type", "application/json")

	json, err := json.Marshal(nome)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(json)

	w.WriteHeader(http.StatusOK)

}
func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	cs, err := rt.AuthenticationApi(r)

	if err != nil {
		http.Error(w, "error: authentication user createChat"+err.Error(), http.StatusUnauthorized)
		return
	}
	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "error parsing multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("propic")
	var photo []byte
	if err != nil {

		noPhotoPath := filepath.Join("webui", "src", "assets", "NoPhoto.png")
		photo, err = os.ReadFile(noPhotoPath)
		if err != nil {
			http.Error(w, "error reading default photo: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		defer file.Close()

		photo, err = validatePhoto(file, handler, err)
		if err != nil {
			http.Error(w, "error: bad input"+err.Error(), http.StatusBadRequest)
			return
		}
	}

	groupName := r.FormValue("nome_chat")
	membriStr := r.FormValue("membri")

	var membri []int
	err = json.Unmarshal([]byte(membriStr), &membri)
	if err != nil {
		http.Error(w, "error decoding membri: "+err.Error(), http.StatusBadRequest)
		return
	}

	// if err := json.NewDecoder(r.Body).Decode(&richiesta); err != nil {
	// 	http.Error(w, "error unmarshal "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	if len(groupName) == 0 {
		http.Error(w, "createGroup error nome chat troppo corto, inserire almeno un carattere:", http.StatusBadRequest)
		return

	}

	id_group, err := rt.db.CreateGroup(cs, groupName, photo, membri)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id_group": id_group,
	}
	w.Header().Set("content-type", "application/json")
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(json)

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
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var name Group
	cs, err := rt.AuthenticationApi(r)

	if err != nil {
		http.Error(w, "error: authentication user createChat"+err.Error(), http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(name.NomeChat) < 1 {
		http.Error(w, "group name too short, at least 1 character", http.StatusBadRequest)
		return
	}
	if len(name.NomeChat) > 16 {
		http.Error(w, "group name too long, up to 16 character", http.StatusBadRequest)
		return
	}

	err = rt.db.SetGroupName(cs, name.IdChat, name.NomeChat)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cs, err := rt.AuthenticationApi(r)

	if err != nil {
		http.Error(w, "error: authentication user createChat"+err.Error(), http.StatusUnauthorized)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "error parsing multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	photo_multipart, handler, err := r.FormFile("propic")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	photo, err := validatePhoto(photo_multipart, handler, err)
	if err != nil {
		http.Error(w, "error: bad input"+err.Error(), http.StatusBadRequest)
		return
	}
	idChat_tmp := r.FormValue("id_chat")
	idChat, err := strconv.Atoi(idChat_tmp)
	if err != nil {
		http.Error(w, "error setGroupPhoto: conversion idChat_tmp (string) to idChat (int)", http.StatusBadRequest)
	}

	fmt.Println(idChat)

	err = rt.db.SetGroupPhoto(cs, idChat, photo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var gruppo Group
	cs, err := rt.AuthenticationApi(r)

	if err != nil {
		http.Error(w, "error: authentication user createChat"+err.Error(), http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&gruppo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = rt.db.AddToGroup(cs, gruppo.IdChat, gruppo.Membri)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
