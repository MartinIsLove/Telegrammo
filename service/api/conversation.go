package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) createChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var richiesta Chat
	cs := r.Header.Get("cs") // cs stà per check session
	if cs == "" {
		http.Error(w, "non è stato restituito alcun autenticatore", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&richiesta.Id)
	if err != nil {
		http.Error(w, "error: conversion json to username"+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(cs)

	if err != nil {
		http.Error(w, "error: conversion string to integer  "+err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}
