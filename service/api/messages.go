package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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

	// if err := json.NewDecoder(r.Body).Decode(&richiesta); err != nil {
	// 	http.Error(w, "sendMessage:"+err.Error(), http.StatusBadRequest)
	// 	return
	// }
	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "error parsing multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("photo")
	var photo []byte
	if err != nil {

		photo = nil
	} else {
		defer file.Close()
		photo, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "error reading file: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	id_chat_tmp := r.FormValue("id_chat")
	testo := r.FormValue("testo")

	id_chat, err := strconv.Atoi(id_chat_tmp)
	if err != nil {
		http.Error(w, "error: conversion id_chat_tmp (string) to id_chat (int): "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(testo) < 1 && photo == nil {
		http.Error(w, "messaggio troppo corto", http.StatusBadRequest)
		return
	}

	erro := rt.db.SendMessage(cs, id_chat, testo, photo)
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
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var emoji Comment

	cs, err := rt.AuthenticationApi(r)
	if err != nil {
		http.Error(w, "error: authentication user sendMessage"+err.Error(), http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&emoji); err != nil {
		http.Error(w, "commentMessage:"+err.Error(), http.StatusBadRequest)
		return
	}

	err = rt.db.CommentMessage(cs, emoji.Id_mes, emoji.Emoji, emoji.Id_chat)
	if err != nil {
		http.Error(w, "commentMessage:"+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
