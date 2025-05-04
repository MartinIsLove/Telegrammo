package api

import (
	"encoding/json"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, "error in database: "+err.Error(), http.StatusInternalServerError)
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

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(json)

}
func (rt *_router) checkChatNames(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var richiesta []ChatUtente // cerca tutte le chat che iniziano per una data stringa
	var str string
	cs, err := rt.AuthenticationApi(r)

	if err != nil {
		http.Error(w, "error: authentication user checknames "+err.Error(), http.StatusUnauthorized)
		return
	}
	str = ps.ByName("nomeChat")
	if len(str) == 0 {
		http.Error(w, "error: username too short", http.StatusBadRequest)
		return
	}

	if len(str) > 16 {
		http.Error(w, "too long name", http.StatusBadRequest)
		return
	}
	if str == "$" {
		w.WriteHeader(http.StatusOK)
		return
	}
	tmp, err := rt.db.CheckChatNames(cs, str)

	if err != nil && strings.HasPrefix(err.Error(), "error in authentication Checknames:") {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && strings.HasPrefix(err.Error(), "CheckChatNames: no chats or groups find:") {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error in database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	richiesta = make([]ChatUtente, len(tmp))
	for i, chat := range tmp {
		richiesta[i] = NewChatUtente(chat)
	}

	w.Header().Set("content-type", "application/json")
	json, err := json.Marshal(richiesta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(json)

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
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
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

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(json)

}
func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	cs, err := rt.AuthenticationApi(r)
	var messaggi []Mess

	if err != nil {
		http.Error(w, "error: authentication user getConversation "+err.Error(), http.StatusUnauthorized)
		return
	}

	id_chat_tmp := ps.ByName("idChat")
	id_chat, err := strconv.Atoi(id_chat_tmp)
	if err != nil {
		http.Error(w, "error getConversation: conversion id_chat_tmp (string) to id_chat (int)", http.StatusBadRequest)
		return
	}

	gp, str, tmp, err := rt.db.GetConversation(cs, id_chat)
	if err != nil && strings.HasPrefix(err.Error(), "error in authentication GetConversation:") {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && (strings.HasPrefix(err.Error(), "GetConversation: no messages found:") || strings.HasPrefix(err.Error(), "GetConversation: no membri found:")) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil && (strings.HasPrefix(err.Error(), "GetConversation: you don't belong to this group")) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
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

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(json)

}
func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	cs, err := rt.AuthenticationApi(r)

	if err != nil {
		http.Error(w, "error: authentication user createGroup "+err.Error(), http.StatusUnauthorized)
		return
	}
	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "error parsing multipart form: "+err.Error(), http.StatusInternalServerError)
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
	// membriStr := r.FormValue("membri")

	// var membri []int
	// err = json.Unmarshal([]byte(membriStr), &membri)
	// if err != nil {
	// 	http.Error(w, "error decoding membri: "+err.Error(), http.StatusBadRequest)
	// 	return
	// }
	membriValues := r.PostForm["membri"]
	if len(membriValues) == 0 {
		http.Error(w, "membri field is required", http.StatusBadRequest)
		return
	}

	membri := make([]int, 0, len(membriValues))
	for _, v := range membriValues {
		id, err := strconv.Atoi(v)
		if err != nil {
			http.Error(w, "invalid member id: "+v, http.StatusBadRequest)
			return
		}
		membri = append(membri, id)
	}

	if len(groupName) == 0 {
		http.Error(w, "createGroup error chat name too short, at least one character:", http.StatusBadRequest)
		return

	}

	id_group, err := rt.db.CreateGroup(cs, groupName, photo, membri)

	// nel controllo degli errori non necessito di tornare un errore nel caso in cui l'utente scelto da aggiungere al gruppo non esista, poiche in automatico aggiunge solo quelli che esistono
	if err != nil && strings.HasPrefix(err.Error(), "error in authentication CreateGroup:") {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
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

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(json)

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
	if err != nil && strings.HasPrefix(err.Error(), "error in authentication LeaveGroup:") {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && (strings.HasPrefix(err.Error(), "LeaveGroup: you don't belong to this group") || strings.HasPrefix(err.Error(), "LeaveGroup: this isn't a group:")) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	cs, err := rt.AuthenticationApi(r)

	if err != nil {
		http.Error(w, "error setGroupName: authentication user createChat "+err.Error(), http.StatusUnauthorized)
		return
	}

	var request struct {
		IdChat   int    `json:"id_chat"`
		NomeChat string `json:"nome_chat"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "error setGroupName: body decode error"+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(request.NomeChat) < 1 {
		http.Error(w, "error setGroupName: group name too short, at least 1 character", http.StatusBadRequest)
		return
	}
	if len(request.NomeChat) > 16 {
		http.Error(w, "error setGroupName: group name too long, up to 16 character", http.StatusBadRequest)
		return
	}

	err = rt.db.SetGroupName(cs, request.IdChat, request.NomeChat)

	if err != nil && strings.HasPrefix(err.Error(), "SetGroupName: error in authentication:") {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && strings.HasPrefix(err.Error(), "SetGroupName: error you can't change a name of a group that you don't partecipate:") {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cs, err := rt.AuthenticationApi(r)

	if err != nil {
		http.Error(w, "error setGroupPhoto: authentication user createChat"+err.Error(), http.StatusUnauthorized)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "error setGroupPhoto: parsing multipart form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	photo_multipart, handler, err := r.FormFile("propic")
	if err != nil {
		http.Error(w, "setGroupPhoto error: form file error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	photo, err := validatePhoto(photo_multipart, handler, err)
	if err != nil {
		http.Error(w, "setGroupPhoto: error bad input photo(the photo must be 1024*1024)"+err.Error(), http.StatusBadRequest)
		return
	}
	idChat_tmp := r.FormValue("id_chat")
	idChat, err := strconv.Atoi(idChat_tmp)
	if err != nil {
		http.Error(w, "error setGroupPhoto: conversion idChat_tmp (string) to idChat (int)", http.StatusBadRequest)
	}

	err = rt.db.SetGroupPhoto(cs, idChat, photo)
	if err != nil && strings.HasPrefix(err.Error(), "SetGroupPhoto: error in authentication:") {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && strings.HasPrefix(err.Error(), "SetGroupPhoto: error you can't change a name of a group that you don't partecipate:") {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var gruppo Group
	cs, err := rt.AuthenticationApi(r)

	if err != nil {
		http.Error(w, "error addToGroup: authentication user "+err.Error(), http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&gruppo); err != nil {
		http.Error(w, "error addToGroup: decoding from body"+err.Error(), http.StatusBadRequest)
		return
	}

	err = rt.db.AddToGroup(cs, gruppo.IdChat, gruppo.Membri)
	if err != nil && strings.HasPrefix(err.Error(), "AddToGroup: error in authentication:") {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && strings.HasPrefix(err.Error(), "AddToGroup: error you can't add user to a group that you don't partecipate:") {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
