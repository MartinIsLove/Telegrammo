package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Richiesta struct {
	Username string `json:"username"`
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var richiestaLogin Richiesta

	err := json.NewDecoder(r.Body).Decode(&richiestaLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// fmt.Println(richiestaLogin)

}
