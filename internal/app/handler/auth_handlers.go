package handler

import (
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-session/session"
	"net/http"
	"net/url"
)

func UserAuthorizationHandler(w http.ResponseWriter, r *http.Request) (string, error) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return "", err
	}

	uid, ok := store.Get("user_id")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}
		store.Set("return_uri", r.Form)
		store.Save()

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return "", nil
	}

	userId := uid.(string)
	store.Delete("user_id")
	store.Save()
	return userId, nil
}

func AuthorizeHandler(srv *server.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, err := session.Start(r.Context(), w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var form url.Values
		if v, ok := store.Get("return_uri"); ok {
			form = v.(url.Values)
		}
		r.Form = form

		store.Delete("return_uri")
		store.Save()

		err = srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func TokenHandler(srv *server.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
