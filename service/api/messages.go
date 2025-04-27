package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	cs, err := rt.AuthenticationApi(r)
	if err != nil {
		http.Error(w, "error authentication user sendMessage:"+err.Error(), http.StatusUnauthorized)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "error parsing multipart form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("photo")
	var photo []byte
	if err != nil {

		photo = nil
	} else {
		defer file.Close()
		photo, err = validatePhoto(file, handler, err)
		//photo, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "error reading file(the photo must be 1024*1024): "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	id_chat_tmp := r.FormValue("id_chat")
	testo := r.FormValue("testo")
	id_reply_tmp := r.FormValue("id_reply")

	id_reply, err := strconv.Atoi(id_reply_tmp)
	if err != nil {
		http.Error(w, "error: conversion id_forward_tmp (string) to id_forward (int): "+err.Error(), http.StatusBadRequest)
		return
	}

	id_chat, err := strconv.Atoi(id_chat_tmp)
	if err != nil {
		http.Error(w, "error: conversion id_chat_tmp (string) to id_chat (int): "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(testo) < 1 && photo == nil {
		http.Error(w, "messaggio troppo corto", http.StatusBadRequest)
		return
	}

	erro := rt.db.SendMessage(cs, id_chat, testo, photo, id_reply)
	if erro != nil && strings.HasPrefix(erro.Error(), "error in authentication SendMessage:") {
		http.Error(w, erro.Error(), http.StatusUnauthorized)
		return
	}
	if erro != nil {
		http.Error(w, erro.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var richiesta RequestData
	cs, err := rt.AuthenticationApi(r)
	if err != nil {
		http.Error(w, "error authentication user forwardMessage:"+err.Error(), http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&richiesta); err != nil {
		http.Error(w, "forwardMessage:"+err.Error(), http.StatusBadRequest)
		return
	}

	err = rt.db.ForwardMessage(cs, richiesta.IdChat, richiesta.IdMes, richiesta.IdForward)
	if err != nil && (strings.HasPrefix(err.Error(), "error in authentication ForwardMessage:") || strings.HasPrefix(err.Error(), "ForwardMessage: error you don't belong to the chat")) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "error database ForwardMessage: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var emoji Comment

	cs, err := rt.AuthenticationApi(r)
	if err != nil {
		http.Error(w, "error authentication user commentMessage: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&emoji); err != nil {
		http.Error(w, "commentMessage:"+err.Error(), http.StatusBadRequest)
		return
	}

	err = rt.db.CommentMessage(cs, emoji.Id_mes, emoji.Emoji, emoji.Id_chat)
	if err != nil && (strings.HasPrefix(err.Error(), "error in authentication CommentMessage:") || strings.HasPrefix(err.Error(), "CommentMessage: you don't belong to the chat")) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "commentMessage:"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var emoji Comment

	cs, err := rt.AuthenticationApi(r)
	if err != nil {
		http.Error(w, "error authentication user uncommentMessage:"+err.Error(), http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&emoji); err != nil {
		http.Error(w, "uncommentMessage:"+err.Error(), http.StatusBadRequest)
		return
	}

	err = rt.db.UncommentMessage(cs, emoji.Id_mes, emoji.Id_chat)
	if err != nil && (strings.HasPrefix(err.Error(), "error in authentication UncommentMessage:") || strings.HasPrefix(err.Error(), "UncommentMessage: you don't belong to the chat")) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "uncommentMessage:"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var messaggio Mess

	cs, err := rt.AuthenticationApi(r)
	if err != nil {
		http.Error(w, "error authentication user deleteMessage "+err.Error(), http.StatusUnauthorized)
		return
	}

	id_mes_tmp := ps.ByName("idMes")
	id_mes, err := strconv.Atoi(id_mes_tmp)
	if err != nil {
		http.Error(w, "error deleteMessage: conversion id_mes_tmp (string) to id_mes (int)", http.StatusBadRequest)
	}

	if err := json.NewDecoder(r.Body).Decode(&messaggio); err != nil {
		http.Error(w, "uncommentMessage:"+err.Error(), http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteMessage(cs, id_mes, messaggio.IdForward, messaggio.IdChat)
	if err != nil && (strings.HasPrefix(err.Error(), "error in authentication DeleteMessage:")) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && (strings.HasPrefix(err.Error(), "DeleteMessage: error database DELETE not successful message don't find")) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "deleteMessage:"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
