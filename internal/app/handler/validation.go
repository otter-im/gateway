package handler

import (
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/server"
	"net/http"
)

func validationBearerToken(server *server.Server, w http.ResponseWriter, r *http.Request) (oauth2.TokenInfo, error) {
	token, err := server.ValidationBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return token, nil
}
