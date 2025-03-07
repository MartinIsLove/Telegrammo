package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	cs := r.Header.Get("cs") // cs stà per check session
	if cs == "" {
		http.Error(w, "non è stato restituito alcun autenticatore", http.StatusBadRequest)
	}
}
