package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var richiesta Utente
	cs := r.Header.Get("cs") // cs stà per check session
	if cs == "" {
		http.Error(w, "non è stato restituito alcun autenticatore", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&richiesta)
	if err != nil {
		http.Error(w, "error: conversion json to username"+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(cs)

	if err != nil {
		http.Error(w, "error: conversion string to integer  "+err.Error(), http.StatusBadRequest)
		return
	}

	err = rt.db.SetMyUserName(richiesta.Username, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
