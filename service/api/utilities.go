package api

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

func (rt *_router) AuthenticationApi(r *http.Request) (int, error) {
	cs := r.Header.Get("cs") // cs stà per check session
	if cs == "" {

		return -1, fmt.Errorf("non è stato restituito alcun autenticatore")
	}
	id, err := strconv.Atoi(cs)

	if err != nil {
		return -1, fmt.Errorf("error in conversion string to int in authenticationApi")
	}
	return id, nil
}
func validatePhoto(photo_multipart multipart.File, handler *multipart.FileHeader, err error) ([]byte, error) {
	var photo []byte

	if err != nil {
		return photo, err
	}

	if handler.Size > 1024*1024 {
		return photo, err
	}

	var erro error
	photo, erro = io.ReadAll(photo_multipart)
	if erro != nil {

		return photo, erro
	}
	return photo, nil
}
