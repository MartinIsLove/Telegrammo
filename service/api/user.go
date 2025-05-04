package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var richiesta Utente

	id, err := rt.AuthenticationApi(r)

	if err != nil {
		http.Error(w, "error: authentication user"+err.Error(), http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&richiesta); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = rt.db.SetMyUserName(richiesta.Username, id)

	if err != nil && strings.HasPrefix(err.Error(), "SetMyUserName: error in authentication:") {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && strings.HasPrefix(err.Error(), "SetMyUserName: username already used, choose another one:") {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var richiesta Utente
	auth, err := rt.AuthenticationApi(r)

	if err != nil {
		http.Error(w, "error: authentication user"+err.Error(), http.StatusUnauthorized)
		return
	}
	if auth > 0 {
		err = r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			http.Error(w, "error parsing multipart form: "+err.Error(), http.StatusInternalServerError)
			return
		}

		photo_multipart, handler, err := r.FormFile("photo")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		photo, err := validatePhoto(photo_multipart, handler, err)
		if err != nil {
			http.Error(w, "setMyPhoto: error bad input(the photo must be 1024*1024)"+err.Error(), http.StatusBadRequest)
			return
		}

		richiesta.Id = auth
		richiesta.Propic = photo
		err = rt.db.SetMyPhoto(richiesta.Propic, richiesta.Id)

		if err != nil && strings.HasPrefix(err.Error(), "SetMyPhoto: error in authentication:") {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	w.WriteHeader(http.StatusNoContent)
}
func (rt *_router) getMyUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var richiesta Utente
	auth, err := rt.AuthenticationApi(r)

	if err != nil {
		http.Error(w, "error: authentication user"+err.Error(), http.StatusUnauthorized)
		return
	}
	if auth > 0 {
		var string_id string = ps.ByName("id")
		id, err := strconv.Atoi(string_id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		richiesta.Username, richiesta.Propic, richiesta.Id, err = rt.db.GetMyUser(id)

		if err != nil && strings.HasPrefix(err.Error(), "GetMyUser: error in authentication:") {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
	if err != nil && strings.HasPrefix(err.Error(), "CheckNames: no user found:") {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error in database: "+err.Error(), http.StatusInternalServerError)
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

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(json)

}
